package cache

import (
	"github.com/joaovrivero/rinha-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.CacheURL,
		PoolSize: 100,
	})

	return nil
}
