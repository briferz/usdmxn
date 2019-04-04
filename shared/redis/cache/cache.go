package cache

import (
	"github.com/briferz/usdmxn/shared/redis/keys"
	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func (c *Cache) Set(key string, data []byte) error {
	cmd := c.client.Set(keys.CacheKey(key), data, 0)
	return cmd.Err()
}

func (c *Cache) Get(key string) ([]byte, error) {
	cmd := c.client.Get(keys.CacheKey(key))
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	return cmd.Bytes()
}
