package env

import (
	"fmt"
	"os"
)

const envRedisHost = "REDIS_HOST"
const defaultRedisHost = "localhost"
const envRedisPort = "REDIS_PORT"
const defaultRedisPort = "6379"
const envRedisPass = "REDIS_PASS"

func RedisHost() string {
	if host, ok := os.LookupEnv(envRedisHost); ok {
		return host
	} else {
		return defaultRedisHost
	}
}

func RedisPort() string {
	if port, ok := os.LookupEnv(envRedisPort); ok {
		return port
	} else {
		return defaultRedisPort
	}
}

func RedisAddr() string {
	return fmt.Sprintf("%s:%s", RedisHost(), RedisPort())
}

func RedisPass()string{
	return os.Getenv(envRedisPass)
}
