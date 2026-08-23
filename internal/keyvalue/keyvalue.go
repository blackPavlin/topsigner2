package keyvalue

import (
	"github.com/redis/go-redis/v9"

	"github.com/bboykiv/topsigner/internal/config"
)

func NewClient(config *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        config.Redis.Addr,
		Password:    config.Redis.Password,
		DB:          config.Redis.DB,
		DialTimeout: config.Redis.DialTimeout,
	})
}
