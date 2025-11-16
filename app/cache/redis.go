package cache

import (
	"context"
	"messaging-system/config"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client      *redis.Client
	cacheConfig *config.Redis
}

func NewRedis(conf *config.Redis) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.URI,
		Password: conf.Password,
		DB:       conf.DB,
	})

	return &Cache{
		client:      client,
		cacheConfig: conf,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	return c.client.Set(ctx, key, value, c.cacheConfig.TTL).Err()
}
