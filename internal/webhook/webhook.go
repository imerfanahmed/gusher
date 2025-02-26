package webhook

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// WebhookConfig represents a webhook configuration
type WebhookConfig struct {
	AppID    string
	Event    string
	URL      string
	APIToken string
}

var webhookConfigs = make(map[string][]WebhookConfig)

// LoadWebhookConfigs loads webhook configurations from the database
func LoadWebhookConfigs(db *sql.DB) {
	rows, err := db.Query("SELECT app_id, event, url, api_token FROM webhooks")
	if err != nil {
		log.Fatal("Error querying webhooks: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var config WebhookConfig
		if err := rows.Scan(&config.AppID, &config.Event, &config.URL, &config.APIToken); err != nil {
			log.Fatal("Error scanning webhook config: ", err)
		}
		webhookConfigs[config.AppID] = append(webhookConfigs[config.AppID], config)
	}
	log.Printf("Loaded webhook configurations for %d apps", len(webhookConfigs))
}

// TriggerWebhook sends HTTP requests to registered webhooks
func TriggerWebhook(appID, event, channel string, data map[string]interface{}) {
	configs, ok := webhookConfigs[appID]
	if !ok {
		return
	}
	for _, config := range configs {
		if config.Event == event {
			payload := map[string]interface{}{
				"event":   event,
				"channel": channel,
				"data":    data,
			}
			jsonPayload, _ := json.Marshal(payload)
			req, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(jsonPayload))
			if err != nil {
				log.Printf("Webhook request creation failed: %v", err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			if config.APIToken != "" {
				req.Header.Set("Authorization", "Bearer "+config.APIToken)
			}
			go func() {
				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Webhook request failed: %v", err)
				} else {
					resp.Body.Close()
				}
			}()
		}
	}
}