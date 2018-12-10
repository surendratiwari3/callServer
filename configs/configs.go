package configs

import (
	"strings"
	"github.com/spf13/viper"
)
type cacheConfig struct {
	Host         string // CACHE_HOST
	PoolSize     int    // CACHE_POOLSIZE
}

// Config - configuration object
type Config struct {
	Cache          cacheConfig
}

var conf *Config

// GetConfig - Function to get Config
func GetConfig() *Config {
	if conf != nil {
		return conf
	}
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cacheConf := cacheConfig{
		Host:         v.GetString("cache.host"),
		PoolSize:     v.GetInt("cache.poolsize"),
	}

	conf = &Config{
		Cache:     cacheConf,
	}
	return conf
}