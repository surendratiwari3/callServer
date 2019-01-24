package repository

import (
	"time"
	"github.com/go-redis/redis"
	"callServer/configs"
	"callServer/adapters"
)

type redisAdapterRepository struct {
	config    *configs.Config
	cacheConn redis.Cmdable
}


// newCacheConnection - Initializes cache connection
func newRedisConnection(config *configs.Config) (redis.Cmdable, error) {
	redisConn := redis.NewClient(&redis.Options{
		Addr:        config.Cache.Host,
		Password:    "",
		DB:          0,
		ReadTimeout: time.Second,
		PoolSize:    config.Cache.PoolSize,
	})
	_, err := redisConn.Ping().Result()
	return redisConn, err
}

// NewCacheAdapterRepository - Repository layer for cache
func NewRedisAdapterRepository(config *configs.Config) (adapters.RedisAdapter, error) {
	redisConn, err := newRedisConnection(config)
	return &redisAdapterRepository{
		config:    config,
		cacheConn: redisConn,
	}, err
}

//Get - Get value from redis
func (c *redisAdapterRepository) Get(key string) (string, error) {
	data, err := c.cacheConn.Get(key).Result()
	c.cacheConn.Del(key)
	return data, err
}

//Set - Set value to redis
func (c *redisAdapterRepository) Set(key string, value string, duration time.Duration) (string, error) {
	result, err := c.cacheConn.Set(key, value, duration).Result()
	return result, err
}