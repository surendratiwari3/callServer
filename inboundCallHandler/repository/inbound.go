package repository

import (
	esl "github.com/cgrates/fsock"
	"callServer/configs"
	coreUtils "callServer/coreUtils/repository"
	"log/syslog"
	"fmt"
	"time"
	"strings"
	"errors"
	"github.com/beevik/etree"
	"callServer/xmlReader/repository"
	"go/doc"
)

type eslAdapterRepository struct {
	config  *configs.Config
	eslConn *ESLsessions
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
		SenderPools: make(map[string]*esl.FSockPool),
	}
	return
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleHeartbeat(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleChannelAnswer(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
}

// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleChannelPark(eventStr, connId string) {
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	onCallURL := "https://gist.githubusercontent.com/surendratiwari3" +
		"/b5d40e8fdc5e6d3a51bca1b4facecfa9/raw/5e7e83eb72252771b5ffe7145afa89e8c7c10ca2/users.xml"
	xmlDocument := repository.GetDocument(onCallURL)
	responseTag := xmlDocument.SelectElement("Response")
	if responseTag != nil {
		fmt.Println("ROOT element:", responseTag.Tag)
		for _, childResponse := range responseTag.ChildElements() {
			fmt.Println("Element:", childResponse.Tag)
			switch childResponse.Tag {
			case "Dial":
				fmt.Println("Value:", childResponse.Text())
				for _, attr := range childResponse.Attr {
					fmt.Printf("Attribute: %s=%s\n", attr.Key, attr.Value)
				}
			case "Play":
				fmt.Println("I'm an int")
				aCallUUID := eventMap["variable_call_uuid"]
				eslCmd := fmt.Sprintf("uuid_broadcast %s %s aleg", aCallUUID,  childResponse.Text())
				eslPool.SendApiCmd(eslCmd)
			default:
				fmt.Printf("Not Valid Tag " + childResponse.Tag)
			}
		}
	}
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
	didNumber := eventMap["variable_sip_req_user"]
	toNumber := "+442078557350"
	dtmfNumber := eventMap["Caller-Caller-ID-Number"]
	isTollFree := eventMap["variable_telemo_tollfree"]
	trunkIP := "mytest11.pstn.sg1.twilio.com"
	aCallUUID := eventMap["variable_call_uuid"]
	//get the solution type associated with this perticular did-number
	// also get the associated organization and userdetails
	// set those details into channel variables so it can be populated into cdr
	// also get the trunk details from the same query only
	// also get the file to be played
	// also in dialplan set call-type inbound
	// also in dialplan need to set the rules from where we are going to get the calls
	if isTollFree == "true" {
		eslCmd := fmt.Sprintf("originate %s %s",
			"{aled_uuid="+aCallUUID+",dtmf_digits="+dtmfNumber+",callbackbridge=true,origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}[send_dtmf=true]sofia/internal/"+toNumber+"@"+trunkIP,
			"&park()")
		eslPool.SendApiCmd(eslCmd)
	}
}


// Formats the event as map and prints it out
func (eslPool *ESLsessions) handleChannelHangup(eventStr, connId string) {
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
		eslPool.SendApiCmd(eslCmd)
	}
	//response, err := c.eslConn.SendCmd(eslCmd)
}

func (eslPool *ESLsessions) handleChannelDTMF(eventStr, connId string) {
	//need to get the solution type from channel name
	// if callflowtype=CLI_ON_RECV_DTMF then only below call flow will applied
	// Format the event from string into Go's map type
	eventMap := esl.FSEventStrToMap(eventStr, []string{})
	aCallUUID := eventMap["Channel-Call-UUID"]
	getDtmdSendDigits := fmt.Sprintf("uuid_getvar %s dtmf_digits", aCallUUID)
	dtmfDigits, err := eslPool.SendApiCmd(getDtmdSendDigits)
	if err!=nil{

	}
	getbCallUUID := fmt.Sprintf("uuid_getvar %s aled_uuid", aCallUUID)
	bCallUUID, err := eslPool.SendApiCmd(getbCallUUID)
	if err!=nil{

	}
	getsend_dtmf := fmt.Sprintf("uuid_getvar %s send_dtmf", aCallUUID)
	send_dtmf, err := eslPool.SendApiCmd(getsend_dtmf)
	if err!=nil{

	}
	dtmfDigitrecv := eventMap["DTMF-Digit"]
	answerState := eventMap["Answer-State"]
	if (dtmfDigitrecv == "1" && answerState == "answered" && send_dtmf == "true") {
		eslCmd := fmt.Sprintf("uuid_send_dtmf %s %s@150", aCallUUID, dtmfDigits)
		eslPool.SendApiCmd(eslCmd)
		eslCmd = fmt.Sprintf("uuid_setvar %s send_dtmf", aCallUUID)
		eslPool.SendApiCmd(eslCmd)
		eslCmd = fmt.Sprintf("uuid_bridge %s %s", aCallUUID, bCallUUID)
		eslPool.SendApiCmd(eslCmd)
	}
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
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
	go func() { // Start reading in own goroutine, return on error
		if err := eslClient.ReadEvents(); err != nil {
			errChan <- err
		}
	}()
	if fsSenderPool, err := esl.NewFSockPool(5, fsAddr, config.EslConfig.Password, 1, 10,
		make(map[string][]func(string, string)), make(map[string][]string), l, connectionUUID); err != nil {
		return fmt.Errorf("Cannot connect FreeSWITCH senders pool, error: %s", err.Error())
	} else if fsSenderPool == nil {
		return errors.New("Cannot connect FreeSWITCH senders pool.")
	} else {
		eslPool.SenderPools[connectionUUID] = fsSenderPool
	}
	eslPool.Conns[connectionUUID] = &eslAdapterRepository{
		config:  config,
		eslConn: eslPool,
	}
	err = <-errChan // Will keep the Connect locked until the first error in one of the connections
	return err
}

// NewCacheAdapterRepository - Repository layer for cache
func NewInboundESLRepository(config *configs.Config, eslPool *ESLsessions) (error) {
	err := newESLConnection(config, eslPool)
	return err
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


//
//XML Reader
//Get the XML
/*
<?xml version="1.0" encoding="UTF-8"?>
<Response>
<Play loop="10">https://api.twilio.com/cowbell.mp3</Play>
</Response>
	//https://www.twilio.com/docs/voice/twiml/play
The <Play> verb plays an audio file back to the caller. Twilio retrieves the file from a URL that you provide.
The <Play> verb supports the following attributes that modify its behavior:

Attribute Name	Allowed Values	Default Value
loop	integer >= 0	1
digits	integer >= 0, w	no default digits for Play
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Play digits="wwww3"></Play>
</Response>
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Play>https://api.twilio.com/cowbell.mp3</Play>
</Response>
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Play loop="10">https://api.twilio.com/cowbell.mp3</Play>
</Response>
ANSWER_URL
FALLBACK_URL
HANGUP_URL
<?xml version="1.0" encoding="UTF-8"?>
<Response>
     <Say>Hello World</Say>
</Response>
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Play loop="10">https://api.twilio.com/cowbell.mp3</Play>
</Response>
