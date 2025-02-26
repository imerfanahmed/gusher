package config

import (
	"database/sql"
	"errors"
	"log"

	"github.com/imerfanahmed/gusher/internal/cache"

	"github.com/go-redis/redis/v8"
)

// AppConfig represents an application configuration
type AppConfig struct {
	AppID     string
	AppKey    string
	AppSecret string
}

var ErrAppNotFound = errors.New("app not found")

// LoadAppConfigs loads app configurations into Redis
func LoadAppConfigs(db *sql.DB, redisClient *redis.Client) {
	rows, err := db.Query("SELECT id, `key`, secret FROM apps")
	if err != nil {
		log.Fatal("Error querying apps: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var config AppConfig
		if err := rows.Scan(&config.AppID, &config.AppKey, &config.AppSecret); err != nil {
			log.Fatal("Error scanning app config: ", err)
		}
		if err := cache.StoreInRedis(redisClient, config.AppKey, config); err != nil {
			log.Printf("Failed to store app %s in Redis: %v", config.AppKey, err)
		}
	}
	log.Printf("Loaded app configurations into Redis")
}
