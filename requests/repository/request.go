package repository

import (
	"net/http"
	"github.com/labstack/echo"
	"callServer/adapters"
	"fmt"
	"strings"
	"callServer/xmlReader/repository"
)

// Controller - struct to logically bind all the controller functions
type Controller struct {
	ESLClient adapters.ESLAdapter
}

func NewRequestController(e *echo.Echo, eslAdapter adapters.ESLAdapter) {
	requestHandler := &Controller{
		ESLClient: eslAdapter,
	}
	e.POST("v1/Account/:auth_id/Call/", requestHandler.call)
}

func (a *Controller) call(c echo.Context) error {
	var callDetails map[string]interface{}

	err := c.Bind(&callDetails)
	response := make(map[string]interface{})
	if err != nil {
		response["error"] = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	//authId := c.Param("auth_id")
	//fromNumber := callDetails["to"].(string)
	toNumber := callDetails["to"].(string)
	didNumber := callDetails["from"].(string)
	//trunkIP := "115.248.91.197"
	trunkIP := "mytest11.pstn.twilio.com"
	response["message"] = "Call is Created"

//	originateCommand := fmt.Sprintf("originate %s %s",
//		"{origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/internal/"+toNumber+"@"+trunkIP,
//		"&bridge({origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/external/"+fromNumber+"@"+trunkIP+")")
	originateCommand := fmt.Sprintf("originate %s %s",
              "{origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/internal/"+toNumber+"@"+trunkIP,
              "&park()")
	resp,_ := a.ESLClient.SendBgApiCmd(originateCommand)
	respField := strings.Fields(resp)
	fmt.Println(respField)
	response["call_uuid"] = string(respField[0])
	answerURL := callDetails["answer_url"].(string)
	xmlDocument := repository.GetDocument(answerURL)
	responseTag := xmlDocument.SelectElement("Response")
	if responseTag != nil {
		fmt.Println("ROOT element:", responseTag.Tag)
		for _, childResponse := range responseTag.ChildElements() {
		fmt.Println("Element:", childResponse.Tag)
		switch childResponse.Tag {
			case "Play":
				aCallUUID := response["call_uuid"]
				eslCmd := fmt.Sprintf("uuid_broadcast %s %s aleg", aCallUUID, childResponse.Text())
				a.ESLClient.SendApiCmd(eslCmd)
			default:
				fmt.Printf("Not Valid Tag " + childResponse.Tag)
				}
		}
	}
	return c.JSON(http.StatusOK, response)
}
