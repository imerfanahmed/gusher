package websocket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/imerfanahmed/gusher/internal/cache"
	"github.com/imerfanahmed/gusher/internal/config"
	"github.com/imerfanahmed/gusher/internal/logging"
	"github.com/imerfanahmed/gusher/internal/webhook"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ChannelManager manages WebSocket connections and channel state
type ChannelManager struct {
	channels         map[string][]*websocket.Conn
	subscriptions    map[*websocket.Conn]map[string]struct{}
	channelOccupancy map[string]int
	mu               sync.RWMutex
}

var cm = ChannelManager{
	channels:         make(map[string][]*websocket.Conn),
	subscriptions:    make(map[*websocket.Conn]map[string]struct{}),
	channelOccupancy: make(map[string]int),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handler handles WebSocket connections
func Handler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		cfg, err := getAppConfig(db, redisClient, key)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid app key", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		socketID := generateSocketID(conn)
		logging.LogEvent(cfg.AppID, "", "connection", socketID)

		cm.mu.Lock()
		cm.subscriptions[conn] = make(map[string]struct{})
		cm.mu.Unlock()

		go handleMessages(conn, cfg, socketID)
	}
}

// handleMessages processes incoming messages
func handleMessages(conn *websocket.Conn, cfg config.AppConfig, socketID string) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var event map[string]interface{}
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}
		eventName, _ := event["event"].(string)
		channel, _ := event["channel"].(string)
		data, _ := event["data"].(map[string]interface{})

		switch eventName {
		case "pusher:subscribe":
			cm.Subscribe(conn, channel)
			if cm.IsChannelOccupied(channel) {
				logging.LogEvent(cfg.AppID, channel, "channel_occupied", socketID)
				webhook.TriggerWebhook(cfg.AppID, "channel_occupied", channel, nil)
			}
		case "pusher:unsubscribe":
			cm.Unsubscribe(conn, channel)
			if cm.IsChannelVacated(channel) {
				logging.LogEvent(cfg.AppID, channel, "channel_vacated", socketID)
				webhook.TriggerWebhook(cfg.AppID, "channel_vacated", channel, nil)
			}
		case "client_event":
			cm.Broadcast(channel, eventName, data)
			logging.LogEvent(cfg.AppID, channel, eventName, socketID)
			webhook.TriggerWebhook(cfg.AppID, eventName, channel, data)
		case "member_added", "member_removed":
			logging.LogEvent(cfg.AppID, channel, eventName, socketID)
			webhook.TriggerWebhook(cfg.AppID, eventName, channel, data)
		}
	}
	cm.UnsubscribeAll(conn)
}

// Helper functions
func getAppConfig(db *sql.DB, redisClient *redis.Client, key string) (config.AppConfig, error) {
	cfg, err := cache.FetchFromRedis(redisClient, key)
	if err == nil {
		return cfg, nil
	}
	cfg, err = fetchFromDatabase(db, key)
	if err != nil {
		return config.AppConfig{}, err
	}
	if err := cache.StoreInRedis(redisClient, key, cfg); err != nil {
		log.Printf("Failed to store app %s in Redis: %v", key, err)
	}
	return cfg, nil
}

func fetchFromDatabase(db *sql.DB, key string) (config.AppConfig, error) {
	var cfg config.AppConfig
	err := db.QueryRow("SELECT id, `key`, secret FROM apps WHERE `key` = ?", key).
		Scan(&cfg.AppID, &cfg.AppKey, &cfg.AppSecret)
	if err == sql.ErrNoRows {
		return config.AppConfig{}, config.ErrAppNotFound
	} else if err != nil {
		return config.AppConfig{}, err
	}
	return cfg, nil
}

func generateSocketID(conn *websocket.Conn) string {
	return fmt.Sprintf("%p", conn)
}

// ChannelManager methods
func (cm *ChannelManager) Subscribe(conn *websocket.Conn, channel string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if _, ok := cm.channels[channel]; !ok {
		cm.channels[channel] = []*websocket.Conn{}
	}
	cm.channels[channel] = append(cm.channels[channel], conn)
	cm.subscriptions[conn][channel] = struct{}{}
	cm.channelOccupancy[channel]++
}

func (cm *ChannelManager) Unsubscribe(conn *websocket.Conn, channel string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if channels, ok := cm.subscriptions[conn]; ok {
		delete(channels, channel)
	}
	if conns, ok := cm.channels[channel]; ok {
		for i, c := range conns {
			if c == conn {
				cm.channels[channel] = append(conns[:i], conns[i+1:]...)
				cm.channelOccupancy[channel]--
				break
			}
		}
	}
}

func (cm *ChannelManager) UnsubscribeAll(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if channels, ok := cm.subscriptions[conn]; ok {
		for channel := range channels {
			cm.Unsubscribe(conn, channel)
		}
		delete(cm.subscriptions, conn)
	}
}

func (cm *ChannelManager) Broadcast(channel, event string, data map[string]interface{}) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	message := map[string]interface{}{
		"event": event,
		"data":  data,
	}
	jsonMessage, _ := json.Marshal(message)
	for _, conn := range cm.channels[channel] {
		conn.WriteMessage(websocket.TextMessage, jsonMessage)
	}
}

func (cm *ChannelManager) IsChannelOccupied(channel string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.channelOccupancy[channel] == 1
}

func (cm *ChannelManager) IsChannelVacated(channel string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.channelOccupancy[channel] == 0
}
