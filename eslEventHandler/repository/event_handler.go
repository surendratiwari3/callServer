package repository

import (
	"fmt"
	"callServer/adapters/repository"
	"strings"
	"time"
)


// Formats the event as map and prints it out
func OnPark(eventMap map[string]string,connId string, eslPool *repository.ESLsessions) {
	// Format the event from string into Go's map type
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

		eslPool.Conns[connId].Originate(eslCmd)
		originateCommand := fmt.Sprintf("originate %s %s",
			"{aled_uuid="+aCallUUID+",dtmf_digits="+dtmfNumber+",callbackbridge=true,origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}[send_dtmf=true]sofia/internal/"+toNumber+"@"+trunkIP,
			"&park()")
		eslCmd = fmt.Sprintf("%s", originateCommand)
		eslPool.Conns[connId].Originate(eslCmd)
	}
}

// Formats the event as map and prints it out
func OnChannelHangup(eventMap map[string]string,connId string, eslPool *repository.ESLsessions) {
	time.Sleep(2000 * time.Millisecond)
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
		eslPool.Conns[connId].Originate(eslCmd)
	}
	//response, err := c.eslConn.SendCmd(eslCmd)
}

func OnDTMF(eventMap map[string]string,connId string, eslPool *repository.ESLsessions) {
	aCallUUID := eventMap["Channel-Call-UUID"]
	getDtmdSendDigits := fmt.Sprintf("uuid_getvar %s dtmf_digits", aCallUUID)
	dtmfDigits := eslPool.Conns[connId].GetVar(getDtmdSendDigits)
	getbCallUUID := fmt.Sprintf("uuid_getvar %s aled_uuid", aCallUUID)
	bCallUUID := eslPool.Conns[connId].GetVar(getbCallUUID)
	getsend_dtmf := fmt.Sprintf("uuid_getvar %s send_dtmf", aCallUUID)
	send_dtmf := eslPool.Conns[connId].GetVar(getsend_dtmf)
	dtmfDigitrecv := eventMap["DTMF-Digit"]
	answerState := eventMap["Answer-State"]
	if (dtmfDigitrecv == "1" && answerState == "answered" && send_dtmf == "true") {
		eslCmd := fmt.Sprintf("uuid_send_dtmf %s %s@150", aCallUUID, dtmfDigits)
		eslPool.Conns[connId].Originate(eslCmd)
		setSendDtmf := fmt.Sprintf("uuid_setvar %s send_dtmf", aCallUUID)
		eslPool.Conns[connId].GetVar(setSendDtmf)
		//time.Sleep(2000 * time.Millisecond)
		originateCommand := fmt.Sprintf("uuid_bridge %s %s", aCallUUID, bCallUUID)
		eslCmd = fmt.Sprintf("%s", originateCommand)
		eslPool.Conns[connId].Originate(eslCmd)
	}
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
}

// Formats the event as map and prints it out
func OnAnswer(eventMap map[string]string,connId string, eslPool *repository.ESLsessions) {
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
}

// Formats the event as map and prints it out
func OnHeartBeat(eventMap map[string]string,connId string, eslPool *repository.ESLsessions) {
	fmt.Printf("%v, connId: %s\n", eventMap, connId)
}
