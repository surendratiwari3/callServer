package adapters

// CacheAdapter - Adapter to talk to cache
type RedisAdapter interface {
	Get(key string) (string, error)
}