package repository

import (
	esl "github.com/cgrates/fsock"
	"callServer/configs"
	"callServer/adapters"
	coreUtils	"callServer/coreUtils/repository"
	"log/syslog"
	"fmt"
)

type eslAdapterRepository struct {
	config    *configs.Config
	eslConn *esl.FSock
}

// Formats the event as map and prints it out
func printHeartbeat( eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n",eventMap, connId)
}

// Formats the event as map and prints it out
func printChannelAnswer( eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n",eventMap, connId)
}

// Formats the event as map and prints it out
func printChannelHangup( eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n",eventMap, connId)
}

func newESLConnection(config *configs.Config)(*esl.FSock, error){
	connectionUUID, err := coreUtils.GenUUID()
	if err != nil {
		panic("not able to generate the connection UUID to connect with FreeSWITCH")
	}
	// Init a syslog writter for our test
	l,errLog := syslog.New(syslog.LOG_INFO, "TestFSock")
	if errLog!=nil {
		panic("not able to connect with syslog")
	}
	fsAddr := fmt.Sprintf("%s:%d",config.EslConfig.Host,config.EslConfig.Port)
	// Filters
	evFilters := make(map[string][]string)
	evFilters["Event-Name"] = append(evFilters["Event-Name"], "CHANNEL_ANSWER")
	evFilters["Event-Name"] = append(evFilters["Event-Name"], "CHANNEL_HANGUP_COMPLETE")

	// We are interested in heartbeats, channel_answer, channel_hangup define handler for them
	evHandlers := map[string][]func(string, string){
		"HEARTBEAT":               {printHeartbeat},
		"CHANNEL_ANSWER":          {printChannelAnswer},
		"CHANNEL_HANGUP_COMPLETE": {printChannelHangup},
	}
	eslClient, err := esl.NewFSock(fsAddr, config.EslConfig.Password, config.EslConfig.Timeout, evHandlers, evFilters, l, connectionUUID)
	if err != nil {
		panic("not able to connect with FreeSWITCH")
	}
	return eslClient, nil
}

// NewCacheAdapterRepository - Repository layer for cache
func NewESLAdapterRepository(config *configs.Config) (adapters.ESLAdapter, error) {
	eslClient, err := newESLConnection(config)
	return &eslAdapterRepository{
		config:    config,
		eslConn: eslClient,
	}, err
}

//Get - Get value from redis
func (c *eslAdapterRepository) Originate(eslCommand string) (string, error) {
	eslCmd := fmt.Sprintf("bgapi %s", eslCommand)
	c.eslConn.SendCmd(eslCmd)
	return c.eslConn.SendCmd(eslCmd)
	//data, err := c.cacheConn.Get(key).Result()
	//return data, err
}