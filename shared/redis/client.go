package redis

import (
	"fmt"
	"github.com/briferz/usdmxn/shared/redis/env"
	"github.com/go-redis/redis"
)

func Client() (*redis.Client, error) {
	opt := redis.Options{
		Addr:     env.RedisAddr(),
		Password: env.RedisPass(),
	}
	client := redis.NewClient(&opt)

	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("unable to reach Redis: %s", err)
	}

	return client, nil
}
