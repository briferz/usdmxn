package redislimiter

import (
	"errors"
	"fmt"
	"github.com/briferz/usdmxn/shared/limiter"
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"
)

type redisLimiter struct {
	client           *redis.Client
	cleanupPeriod    time.Duration
	periodAllowances int64
}

func New(client *redis.Client, period time.Duration, allowances int64) (limiter.Interface, <-chan error) {
	if client == nil || period <= 0 || allowances <= 0 {
		log.Panicf("bad param(s) (%v / %v / %v)", client, period, allowances)
	}

	lim := &redisLimiter{
		client:           client,
		cleanupPeriod:    period,
		periodAllowances: allowances,
	}

	errCh := lim.periodicReqCounterCleanup()
	return lim, errCh

}

func (l *redisLimiter) increaseCounter(user string) (int64, error) {
	if user == "" {
		return 0, errors.New("empty user")
	}
	cmd := l.client.HIncrBy(usersReqHash, user, 1)
	if err := cmd.Err(); err != nil {
		return 0, fmt.Errorf("error increasing counter for user '%s': %s", user, err)
	}

	return cmd.Val(), nil
}

func (l *redisLimiter) decreaseCounter(user string) error {
	if user == "" {
		return errors.New("empty user")
	}
	cmd := l.client.HIncrBy(usersReqHash, user, -1)
	return cmd.Err()
}

func (l *redisLimiter) periodicReqCounterCleanup() <-chan error {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		ticker := time.NewTicker(l.cleanupPeriod)
		for range ticker.C {
			cmd := l.client.Del(usersReqHash)
			if err := cmd.Err(); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}

	}()
	wg.Wait()
	return errCh
}

func (l *redisLimiter) Allow(user string) (bool, error) {
	currentCounter, err := l.increaseCounter(user)
	if err != nil {
		return false, err
	}
	if currentCounter > l.periodAllowances {
		decreaseErr := l.decreaseCounter(user)
		if decreaseErr != nil {
			return false, decreaseErr
		}
		return false, nil
	}
	return true, nil
}
