package repository

import (
	esl "github.com/cgrates/fsock"
	"callServer/configs"
	"callServer/adapters"
	coreUtils "callServer/coreUtils/repository"
	"log/syslog"
	"fmt"
	"time"
)

type eslAdapterRepository struct {
	config    *configs.Config
	eslConn *esl.FSock
}

// The freeswitch session manager type holding a buffer for the network connection
// and the active sessions
type ESLsessions struct {
	cfg         *configs.Config
	conns       map[string]*eslAdapterRepository // Keep the list here for connection management purposes
	senderPools map[string]*esl.FSockPool  // Keep sender pools here
}

func NewESLsessions(config *configs.Config) (eslPool *ESLsessions) {
	eslPool = &ESLsessions{
		cfg:         config,
		conns:       make(map[string]*eslAdapterRepository),
	}
	return
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
	time.Sleep(2000 * time.Millisecond)
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n",eventMap, connId)
	didNumber := "919967609476"
	toNumber := "919967609476"
	fromNumber := "+17139187788"
	trunkIP := "10.17.112.21"

	originateCommand := fmt.Sprintf("originate %s %s",
		"{origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/internal/"+toNumber+"@"+trunkIP,
		"&bridge({origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/external/"+fromNumber+"@"+trunkIP+")")
	eslCmd := fmt.Sprintf("bgapi %s", originateCommand)
	eslAdapterRepository.Originate(eslCmd)
	//response, err := c.eslConn.SendCmd(eslCmd)
}

func newESLConnection(config *configs.Config, eslPool)(*esl.FSock, error){
	errChan := make(chan error)
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
		"HEARTBEAT":               {eslPool.printHeartbeat},
		"CHANNEL_ANSWER":          {eslPool.printChannelAnswer},
		"CHANNEL_HANGUP_COMPLETE": {eslPool.printChannelHangup},
	}
	eslClient, err := esl.NewFSock(fsAddr, config.EslConfig.Password, config.EslConfig.Timeout, evHandlers, evFilters, l, connectionUUID)
	if err != nil {
		panic("not able to connect with FreeSWITCH")
	}
	eslPool.conns[connectionUUID] = &eslAdapterRepository{
		config:    config,
		eslConn: eslClient,
	}
	go func() { // Start reading in own goroutine, return on error
		if err := eslPool.conns[connectionUUID].eslConn.ReadEvents(); err != nil {
			errChan <- err
		}
	}()
	return eslClient, nil
}

// NewCacheAdapterRepository - Repository layer for cache
func NewESLAdapterRepository(config *configs.Config,eslPool *ESLsessions) (adapters.ESLAdapter, error) {
	eslClient, err := newESLConnection(config, eslPool)
	return &eslAdapterRepository{
		config:    config,
		eslConn: eslClient,
	}, err
}

//Get - Get value from redis
func (c *eslAdapterRepository) Originate(eslCommand string) (string, error) {
	eslCmd := fmt.Sprintf("bgapi %s", eslCommand)
	return c.eslConn.SendCmd(eslCmd)
	//data, err := c.cacheConn.Get(key).Result()
	//return data, err
}