// internal/config/config.go
package config

import (
    "database/sql"
    "log"
    "github.com/imerfanahmed/gusher/internal/types"
)

// LoadAppConfigs loads app configurations from the database
func LoadAppConfigs(db *sql.DB) ([]types.AppConfig, error) {
    rows, err := db.Query("SELECT id, `key`, secret FROM apps")
    if err != nil {
        log.Printf("Error querying apps: %v", err)
        return nil, err
    }
    defer rows.Close()

    var configs []types.AppConfig
    for rows.Next() {
        var config types.AppConfig
        if err := rows.Scan(&config.AppID, &config.AppKey, &config.AppSecret); err != nil {
            log.Printf("Error scanning app config: %v", err)
            return nil, err
        }
        configs = append(configs, config)
    }
    log.Printf("Loaded app configurations from database")
    return configs, nil
}