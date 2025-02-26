// internal/cache/cache.go
package cache

import (
    "context"
    "time"
    "github.com/imerfanahmed/gusher/internal/types"
    "github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(host, port string) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr: host + ":" + port,
    })
    if _, err := client.Ping(context.Background()).Result(); err != nil {
        return nil, err
    }
    return client, nil
}

// FetchFromRedis retrieves an app config from Redis
func FetchFromRedis(redisClient *redis.Client, key string) (types.AppConfig, error) {
    redisKey := "app:" + key
    result, err := redisClient.HGetAll(context.Background(), redisKey).Result()
    if err != nil || len(result) == 0 {
        return types.AppConfig{}, types.ErrAppNotFound
    }
    return types.AppConfig{
        AppID:     result["app_id"],
        AppKey:    key,
        AppSecret: result["app_secret"],
    }, nil
}

// StoreInRedis stores an app config in Redis with a TTL
func StoreInRedis(redisClient *redis.Client, key string, cfg types.AppConfig) error {
    redisKey := "app:" + key
    err := redisClient.HSet(context.Background(), redisKey, map[string]interface{}{
        "app_id":     cfg.AppID,
        "app_secret": cfg.AppSecret,
    }).Err()
    if err != nil {
        return err
    }
    return redisClient.Expire(context.Background(), redisKey, 5*time.Minute).Err()
}