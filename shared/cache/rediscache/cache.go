package rediscache

import (
	"github.com/briferz/usdmxn/shared/cache"
	"github.com/go-redis/redis"
)

type redisCache struct {
	client *redis.Client
}

func New(client *redis.Client) cache.Interface {
	return &redisCache{
		client: client,
	}
}

func (c *redisCache) Set(key string, data []byte) error {
	cmd := c.client.Set(cacheKey(key), data, 0)
	return cmd.Err()
}

func (c *redisCache) Get(key string) ([]byte, error) {
	cmd := c.client.Get(cacheKey(key))
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	return cmd.Bytes()
}
