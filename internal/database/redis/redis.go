package redis

import (
	"context"
	"fmt"

	"cloud-notes/internal/config"

	"github.com/redis/go-redis/v9"
)

type Redis = redis.Client

var ErrNoRows = redis.Nil

func Connect(ctx context.Context, cfg *config.Redis) (*Redis, error) {
	r := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   cfg.DB,
	})

	err := r.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return r, nil
}

func MustConnect(ctx context.Context, cfg *config.Redis) *Redis {
	r, err := Connect(ctx, cfg)
	if err != nil {
		panic(err)
	}

	return r
}
