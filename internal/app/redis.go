package app

import (
	"context"

	"github.com/go-redis/redis/v9"
)

func initRedisClient(ctx context.Context, cfg configer) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword(),
		DB:       int(cfg.RedisDB()),
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
