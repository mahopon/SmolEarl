package redis

import (
	"context"
	"fmt"
	config "github.com/mahopon/SmolEarl/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	AppConfig = config.AppConfig
)

type Redis struct {
	Client *redis.Client
}

func InitRedis() (*Redis, error) {
	r := &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppConfig.RedisHost, AppConfig.RedisPort),
			Password: AppConfig.RedisPass,
			DB:       0,
		}),
	}

	ctx := context.Background()
	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return r, nil
}

func (r *Redis) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()

	if err != nil {
		return err
	}
	return nil
}

func getErrorCode(err error) string {
	if err == redis.Nil {
		return "not_found"
	}
	return err.Error()
}
