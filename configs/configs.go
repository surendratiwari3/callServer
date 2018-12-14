package configs

import (
	"strings"
	"github.com/spf13/viper"
	"fmt"
)
type cacheConfig struct {
	Host         string // CACHE_HOST
	PoolSize     int    // CACHE_POOLSIZE
}

type logConfig struct {
	LogFile		   string
	LogLevel 		string
}

type httpConfig struct {
	HostPort string
}

type eslConfig struct {
	HostPort string
	Password string
	Reconnect int
}

// Config - configuration object
type Config struct {
	Cache          cacheConfig
	Log		   logConfig
	HttpConfig httpConfig
	EslConfig	eslConfig
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

	fmt.Println(v.Get("cache.host"))
	fmt.Println(v.Get("cache_host"))

	cacheConf := cacheConfig{
		Host:         v.GetString("cache.host"),
		PoolSize:     v.GetInt("cache.poolsize"),
	}

	logConf := logConfig{
		LogFile:         v.GetString("log.file"),
		LogLevel:     v.GetString("log.level"),
	}

	httpConf := httpConfig{
		HostPort:         v.GetString("http.host"),
	}
	eslConf := eslConfig{
		HostPort:         v.GetString("esl.host"),
		Password:         v.GetString("esl.password"),
		Reconnect:         v.GetInt("esl.reconnect"),
	}

	conf = &Config{
		Cache:     cacheConf,
		Log: 	logConf,
		HttpConfig: httpConf,
		EslConfig: eslConf,
	}
	return conf
}