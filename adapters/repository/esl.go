package repository

import (
	esl "github.com/cgrates/fsock"
	"callServer/configs"
	"callServer/adapters"
	coreUtils "callServer/coreUtils/repository"
	"log/syslog"
	"fmt"
	"strings"
	"errors"
)

type eslAdapterRepository struct {
	config  *configs.Config
	eslPool *ESLsessions
}

// The freeswitch session manager type holding a buffer for the network connection
// and the active sessions
type ESLsessions struct {
	Cfg         *configs.Config
	SenderPools map[string]*esl.FSockPool        // Keep sender pools here
}

func NewESLsessions(config *configs.Config) (eslPool *ESLsessions) {
	eslPool = &ESLsessions{
		Cfg:   config,
		SenderPools: make(map[string]*esl.FSockPool),
	}
	return
}

func newESLConnection(config *configs.Config, eslPool *ESLsessions) (error) {
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
	if fsSenderPool, err := esl.NewFSockPool(5, fsAddr, config.EslConfig.Password, 1, 10,
		make(map[string][]func(string, string)), make(map[string][]string), l, connectionUUID); err != nil {
		return fmt.Errorf("Cannot connect FreeSWITCH senders pool, error: %s", err.Error())
	} else if fsSenderPool == nil {
		return errors.New("Cannot connect FreeSWITCH senders pool.")
	} else {
		eslPool.SenderPools[connectionUUID] = fsSenderPool
	}
	err = <-errChan // Will keep the Connect locked until the first error in one of the connections
	return err
}

// NewCacheAdapterRepository - Repository layer for cache
func NewESLAdapterRepository(config *configs.Config, eslPool *ESLsessions) (adapters.ESLAdapter, error) {
	err := newESLConnection(config, eslPool)
	return &eslAdapterRepository{
		config:  config,
		eslPool: eslPool,
	}, err
}

//Get - Get value from redis
func (c *eslAdapterRepository) SendBgApiCmd(eslCommand string) (string, error) {
	eslCmd := fmt.Sprintf("bgapi %s", eslCommand)
	resp, err := c.eslPool.SendBGApiCmd(eslCmd)
	respField := strings.Fields(resp)
	uuid := string(respField[2])
	//data, err := c.cacheConn.Get(key).Result()
	return uuid, err
}

//Get - Get value from redis
func (c *eslAdapterRepository) SendApiCmd(eslCommand string) (string) {
	resp, err := c.eslPool.SendApiCmd(eslCommand)
	if err == nil {
		return resp
	}
	return ""
}

func (eslPool *ESLsessions) SendBGApiCmd(eslCommand string) (response string, err error) {
	l, errLog := syslog.New(syslog.LOG_INFO, "TestFSock")
	if errLog != nil {
		panic("not able to connect with syslog")
	}
	for connId, senderPool := range eslPool.SenderPools {
		fsConn, err := senderPool.PopFSock()
		if err != nil {
			l.Err(fmt.Sprintf("<%s> Error on connection id: %s", err.Error(), connId))
			continue
		}
		response, err = fsConn.SendApiCmd(eslCommand)
		senderPool.PushFSock(fsConn)
		return response, err
	}
	return response, err
}

func (eslPool *ESLsessions) SendApiCmd(eslCommand string) (response string, err error) {
	l, errLog := syslog.New(syslog.LOG_INFO, "TestFSock")
	if errLog != nil {
		panic("not able to connect with syslog")
	}
	for connId, senderPool := range eslPool.SenderPools {
		fsConn, err := senderPool.PopFSock()
		if err != nil {
			l.Err(fmt.Sprintf("<%s> Error on connection id: %s", err.Error(), connId))
			continue
		}
		response, err = fsConn.SendApiCmd(eslCommand)
		senderPool.PushFSock(fsConn)
		return response, err
	}
	return response, err
}