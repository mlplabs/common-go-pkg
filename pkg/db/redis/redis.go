package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Database,
	})
}
