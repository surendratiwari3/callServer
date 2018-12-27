package repository

import (
	esl "github.com/cgrates/fsock"
	"callServer/configs"
	"callServer/adapters"
	coreUtils "callServer/coreUtils/repository"
	"log/syslog"
	"fmt"
	"strings"
	"callServer/eslEventHandler/repository"
)

type eslAdapterRepository struct {
	config  *configs.Config
	eslConn *esl.FSock
}

// The freeswitch session manager type holding a buffer for the network connection
// and the active sessions
type ESLsessions struct {
	Cfg         *configs.Config
	Conns       map[string]*eslAdapterRepository // Keep the list here for connection management purposes
	SenderPools map[string]*esl.FSockPool        // Keep sender pools here
}

func NewESLsessions(config *configs.Config) (eslPool *ESLsessions) {
	eslPool = &ESLsessions{
		Cfg:   config,
		Conns: make(map[string]*eslAdapterRepository),
	}
	return
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleHeartbeat(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	repository.OnHeartBeat(eventMap,connId, eslPool)
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleChannelAnswer(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	repository.OnAnswer(eventMap,connId, eslPool)
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleChannelPark(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	repository.OnPark(eventMap,connId, eslPool)
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleChannelHangup(eventStr, connId string) {
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	repository.OnChannelHangup(eventMap,connId, eslPool)
}

func (eslPool *ESLsessions) handleChannelDTMF(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	repository.OnDTMF(eventMap,connId, eslPool)
}

func newESLConnection(config *configs.Config, eslPool *ESLsessions) (*esl.FSock, error) {
	errChan := make(chan error)
	connectionUUID, err := coreUtils.GenUUID()
	if err != nil {
		panic("not able to generate the connection UUID to connect with FreeSWITCH")
	}
	// Init a syslog writter for our test
	l, errLog := syslog.New(syslog.LOG_INFO, "TestFSock")
	if errLog != nil {
		panic("not able to connect with syslog")
	}
	fsAddr := fmt.Sprintf("%s:%d", config.EslConfig.Host, config.EslConfig.Port)
	// Filters
	evFilters := make(map[string][]string)
	evFilters["Event-Name"] = append(evFilters["Event-Name"], "CHANNEL_ANSWER")
	evFilters["Event-Name"] = append(evFilters["Event-Name"], "CHANNEL_HANGUP_COMPLETE")
	evFilters["Event-Name"] = append(evFilters["Event-Name"], "CHANNEL_PARK")
	evFilters["Event-Name"] = append(evFilters["Event-Name"], "DTMF")

	// We are interested in heartbeats, channel_answer, channel_hangup define handler for them
	evHandlers := map[string][]func(string, string){
		"HEARTBEAT":               {eslPool.handleHeartbeat},
		"CHANNEL_ANSWER":          {eslPool.handleChannelAnswer},
		"CHANNEL_HANGUP_COMPLETE": {eslPool.handleChannelHangup},
		"CHANNEL_PARK":            {eslPool.handleChannelPark},
		"DTMF":                    {eslPool.handleChannelDTMF},
	}
	eslClient, err := esl.NewFSock(fsAddr, config.EslConfig.Password, config.EslConfig.Timeout, evHandlers, evFilters, l, connectionUUID)
	if err != nil {
		panic("not able to connect with FreeSWITCH")
	}
	eslPool.Conns[connectionUUID] = &eslAdapterRepository{
		config:  config,
		eslConn: eslClient,
	}
	go func() { // Start reading in own goroutine, return on error
		if err := eslPool.Conns[connectionUUID].eslConn.ReadEvents(); err != nil {
			errChan <- err
		}
	}()
	return eslClient, nil
}

// NewCacheAdapterRepository - Repository layer for cache
func NewESLAdapterRepository(config *configs.Config, eslPool *ESLsessions) (adapters.ESLAdapter, error) {
	eslClient, err := newESLConnection(config, eslPool)
	return &eslAdapterRepository{
		config:  config,
		eslConn: eslClient,
	}, err
}

//Get - Get value from redis
func (c *eslAdapterRepository) Originate(eslCommand string) (string, error) {
	eslCmd := fmt.Sprintf("bgapi %s", eslCommand)
	resp, err := c.eslConn.SendCmd(eslCmd)
	respField := strings.Fields(resp)
	uuid := string(respField[2])
	//data, err := c.cacheConn.Get(key).Result()
	return uuid, err
}

//Get - Get value from redis
func (c *eslAdapterRepository) GetVar(eslCommand string) (string) {
	resp, err := c.eslConn.SendApiCmd(eslCommand)
	//data, err := c.cacheConn.Get(key).Result()
	if err == nil {
		return resp
	}
	return ""
}
