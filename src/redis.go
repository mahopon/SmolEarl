package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisClient holds the Redis client instance
var RedisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppConfig.RedisHost, AppConfig.RedisPort),
		Password: AppConfig.RedisPass,
		DB:       0, // use default DB
	})

	// Test connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
