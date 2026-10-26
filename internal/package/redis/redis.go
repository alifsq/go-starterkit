package redis

import (
	"context"
	"fmt"
	"go-starterkit/internal/config"
	"time"

	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
)

var RedisSet = wire.NewSet(
	NewRedisClient,
)

func NewRedisClient(cfg *config.Config) (*goredis.Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke Redis di %s: %w", cfg.Addr, err)
	}

	return rdb, nil
}
