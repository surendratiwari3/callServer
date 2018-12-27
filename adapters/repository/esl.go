package repository

import (
	esl "github.com/cgrates/fsock"
	"callServer/configs"
	"callServer/adapters"
	coreUtils "callServer/coreUtils/repository"
	"log/syslog"
	"fmt"
	"time"
	"strings"
)

type eslAdapterRepository struct {
	config  *configs.Config
	eslConn *esl.FSock
}

// The freeswitch session manager type holding a buffer for the network connection
// and the active sessions
type ESLsessions struct {
	cfg         *configs.Config
	conns       map[string]*eslAdapterRepository // Keep the list here for connection management purposes
	senderPools map[string]*esl.FSockPool        // Keep sender pools here
}

func NewESLsessions(config *configs.Config) (eslPool *ESLsessions) {
	eslPool = &ESLsessions{
		cfg:   config,
		conns: make(map[string]*eslAdapterRepository),
	}
	return
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) printHeartbeat(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) printChannelAnswer(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleChannelPark(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
	didNumber := eventMap["variable_sip_req_user"]
	toNumber := "+442078557350"
	dtmfNumber := eventMap["Caller-Caller-ID-Number"]
	isTollFree := eventMap["variable_telemo_tollfree"]
	trunkIP := "mytest11.pstn.sg1.twilio.com"
	aCallUUID := eventMap["variable_call_uuid"]
	if isTollFree == "true" {
		uuidSet := fmt.Sprintf("uuid_broadcast %s %s aleg", aCallUUID, "/usr/local/freeswitch/sounds/bridge.wav")
		eslCmd := fmt.Sprintf("%s", uuidSet)
		eslPool.conns[connId].Originate(eslCmd)
		originateCommand := fmt.Sprintf("originate %s %s",
			"{aled_uuid="+aCallUUID+",dtmf_digits="+dtmfNumber+",callbackbridge=true,origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}[send_dtmf=true]sofia/internal/"+toNumber+"@"+trunkIP,
			"&park()")
		eslCmd = fmt.Sprintf("%s", originateCommand)
		eslPool.conns[connId].Originate(eslCmd)
	}
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) printChannelHangup(eventStr, connId string) {
	time.Sleep(2000 * time.Millisecond)
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	didNumber := eventMap["variable_sip_req_user"]
	toNumber := eventMap["Caller-Caller-ID-Number"]
	fromNumber := "+17139187788"
	trunkIP := "mytest11.pstn.sg1.twilio.com"
	dtmfDigits := strings.TrimPrefix(toNumber, "+")
	hangupApplication := eventMap["variable_current_application_data"]
	if (hangupApplication == "ESL_TERMINATE") {
		originateCommand := fmt.Sprintf("originate %s %s",
			"{ignore_early_media=true,origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}[execute_on_answer='send_dtmf "+dtmfDigits+"']sofia/internal/"+toNumber+"@"+trunkIP,
			"&bridge({origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/external/"+fromNumber+"@"+trunkIP+")")
		eslCmd := fmt.Sprintf("bgapi %s", originateCommand)
		eslPool.conns[connId].Originate(eslCmd)
	}
	//response, err := c.eslConn.SendCmd(eslCmd)
}

func (eslPool *ESLsessions) handleChannelDTMF(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	aCallUUID := eventMap["Channel-Call-UUID"]
	//        bCallUUID := eventMap["variable_aled_uuid"]
	//dtmfDigits := eventMap["variable_dtmf_digits"]
	//dtmfDigits := "919967609476"
	getDtmdSendDigits := fmt.Sprintf("uuid_getvar %s dtmf_digits", aCallUUID)
	dtmfDigits := eslPool.conns[connId].GetVar(getDtmdSendDigits)
	getbCallUUID := fmt.Sprintf("uuid_getvar %s aled_uuid", aCallUUID)
	bCallUUID := eslPool.conns[connId].GetVar(getbCallUUID)
	getsend_dtmf := fmt.Sprintf("uuid_getvar %s send_dtmf", aCallUUID)
	send_dtmf := eslPool.conns[connId].GetVar(getsend_dtmf)

	dtmfDigitrecv := eventMap["DTMF-Digit"]
	answerState := eventMap["Answer-State"]
	if (dtmfDigitrecv == "1" && answerState == "answered" && send_dtmf == "true") {
		eslCmd := fmt.Sprintf("uuid_send_dtmf %s %s@150", aCallUUID, dtmfDigits)
		eslPool.conns[connId].Originate(eslCmd)
		setSendDtmf := fmt.Sprintf("uuid_setvar %s send_dtmf", aCallUUID)
		eslPool.conns[connId].GetVar(setSendDtmf)
		//time.Sleep(2000 * time.Millisecond)
		originateCommand := fmt.Sprintf("uuid_bridge %s %s", aCallUUID, bCallUUID)
		eslCmd = fmt.Sprintf("%s", originateCommand)
		eslPool.conns[connId].Originate(eslCmd)
	}
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
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
		"HEARTBEAT":               {eslPool.printHeartbeat},
		"CHANNEL_ANSWER":          {eslPool.printChannelAnswer},
		"CHANNEL_HANGUP_COMPLETE": {eslPool.printChannelHangup},
		"CHANNEL_PARK":            {eslPool.handleChannelPark},
		"DTMF":                    {eslPool.handleChannelDTMF},
	}
	eslClient, err := esl.NewFSock(fsAddr, config.EslConfig.Password, config.EslConfig.Timeout, evHandlers, evFilters, l, connectionUUID)
	if err != nil {
		panic("not able to connect with FreeSWITCH")
	}
	eslPool.conns[connectionUUID] = &eslAdapterRepository{
		config:  config,
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
