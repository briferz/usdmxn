package redistokencreatorvalidator

import (
	"crypto/sha1"
	"fmt"
	"github.com/briferz/usdmxn/middleware/tokenmiddleware/tokencreatorvalidator"
	"github.com/go-redis/redis"
	"time"
)

const (
	redisKeysSet = "apiKeys"
)

type creatorValidator struct {
	client *redis.Client
}

func New(client *redis.Client) tokencreatorvalidator.CreatorValidator {
	return &creatorValidator{
		client: client,
	}
}

func (cv *creatorValidator) Validate(key string) (bool, error) {
	cmd := cv.client.SIsMember(redisKeysSet, key)
	if err := cmd.Err(); err != nil {
		return false, fmt.Errorf("error querying redis for key %s: %s", key, err)
	}

	return cmd.Result()
}

func (cv *creatorValidator) Create() (string, error) {
	timeNano := fmt.Sprint(time.Now().UnixNano())
	h := sha1.New()
	h.Write([]byte(timeNano))
	token := fmt.Sprintf("%x", h.Sum(nil))

	cmd := cv.client.SAdd(redisKeysSet, token)
	if err := cmd.Err(); err != nil {
		return "", fmt.Errorf("error setting new key in redis: %s", err)
	}

	if cmd.Val() == 0 {
		return "", fmt.Errorf("the just generated key was already set in redis")
	}

	return token, nil
}
