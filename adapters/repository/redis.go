package repository

import (
	"time"
	"github.com/go-redis/redis"
	"callServer/configs"
	"callServer/adapters"
	"github.com/Sirupsen/logrus"
)

type redisAdapterRepository struct {
	config    *configs.Config
	cacheConn redis.Cmdable
}


// newCacheConnection - Initializes cache connection
func newCacheConnection(config *configs.Config, logger *logrus.Logger) redis.Cmdable {
	cacheConn := redis.NewClient(&redis.Options{
		Addr:        config.Cache.Host,
		Password:    "",
		DB:          0,
		ReadTimeout: time.Second,
		PoolSize:    config.Cache.PoolSize,
	})
	if cacheConn == nil {
		logger.WithField("redis_host", config.Cache.Host).Errorf("Can't connect to redis")
	}
	return cacheConn
}

// NewCacheAdapterRepository - Repository layer for cache
func NewCacheAdapterRepository(config *configs.Config, logger *logrus.Logger) adapters.RedisAdapter {
	logger.WithField("redis_host", config.Cache.Host).Info("connecting to redis")
	cacheConn := newCacheConnection(config, logger)
	if cacheConn == nil {
		panic("unable to connect to redis")
	}
	return &redisAdapterRepository{
		config:    config,
		cacheConn: cacheConn,
	}
}

//Get - Get value from redis
func (c *redisAdapterRepository) Get(key string) (string, error) {
	data, err := c.cacheConn.Get(key).Result()
	return data, err
}