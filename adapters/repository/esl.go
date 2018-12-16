package repository

import (
	esl "github.com/0x19/goesl"
	"callServer/configs"
	"callServer/adapters"
)

type eslAdapterRepository struct {
	config    *configs.Config
	eslConn *esl.Client
}

func newESLConnection(config *configs.Config)(*esl.Client, error){
	//ESL
	eslClient, err := esl.NewClient(config.EslConfig.Host, config.EslConfig.Port, config.EslConfig.Password, config.EslConfig.Timeout)
	if err != nil {
		panic("not able to connect with FreeSWITCH")
	}
	return eslClient, nil
}

// NewCacheAdapterRepository - Repository layer for cache
func NewESLAdapterRepository(config *configs.Config) (adapters.ESLAdapter, error) {
	eslClient, err := newESLConnection(config)
	go eslClient.Handle()
	return &eslAdapterRepository{
		config:    config,
		eslConn: eslClient,
	}, err
}

//Get - Get value from redis
func (c *eslAdapterRepository) Originate(eslCommand string) (string, error) {
	c.eslConn.BgApi(eslCommand)
	return "hello",nil
	//data, err := c.cacheConn.Get(key).Result()
	//return data, err
}