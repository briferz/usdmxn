package keys

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

const (
	usersReqHash = "users_req_counts_hash"
)

func IncreaseCounter(user string, client *redis.Client) (int64, error) {
	if user == "" {
		return 0, errors.New("empty user")
	}
	cmd := client.HIncrBy(usersReqHash, user, 1)
	if err := cmd.Err(); err != nil {
		return 0, fmt.Errorf("error increasing counter for user '%s': %s", user, err)
	}

	return cmd.Val(), nil
}

func DecreaseCounter(user string, client *redis.Client) error {
	if user == "" {
		return errors.New("empty user")
	}
	cmd := client.HIncrBy(usersReqHash, user, -1)
	return cmd.Err()
}

func PeriodicReqCounterCleanup(client *redis.Client, duration time.Duration) (<-chan error, error) {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		ticker := time.NewTicker(duration)
		for range ticker.C {
			cmd := client.Del(usersReqHash)
			if err := cmd.Err(); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}

	}()
	wg.Wait()
	return errCh, nil
}

func TestForLimiter(client *redis.Client, limit int64, user string) (bool, error) {
	currentCounter, err := IncreaseCounter(user, client)
	if err != nil {
		return false, err
	}
	if currentCounter > limit {
		decreaseErr := DecreaseCounter(user, client)
		if decreaseErr != nil {
			return false, decreaseErr
		}
		return false, nil
	}
	return true, nil
}
