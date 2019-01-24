package adapters

import "time"

// CacheAdapter - Adapter to talk to cache
type RedisAdapter interface {
	Get(key string) (string, error)
	Set(key string, value string, duration time.Duration) (string, error)
}