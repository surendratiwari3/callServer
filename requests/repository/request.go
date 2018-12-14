package repository

import (
	"net/http"
	"github.com/labstack/echo"
	esl "github.com/0x19/goesl"
	coreUtils "callServer/coreUtils/repository"
	"fmt"
)

// Controller - struct to logically bind all the controller functions
type Controller struct {
	ESLClient *esl.Client
}

func NewRequestController(e *echo.Echo, client *esl.Client) {
	requestHandler := &Controller{
		ESLClient: client,
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

	callUUID, err :=  coreUtils.GenUUID()
	if err != nil{
		panic("not able to generate the connection UUID to connect with FreeSWITCH")
	}

	authId := c.Param("auth_id")
	fromNumber := callDetails["from_number"].(string)
	toNumber := callDetails["toNumber"].(string)
	didNumber := callDetails["didNumber"].(string)

	response["message"] = "Successfully Authenticated " + authId
	eslJobid := a.ESLClient.BgApi(fmt.Sprintf("originate %s %s",
		"{origination_uuid="+callUUID+",origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/internal/"+toNumber+"@10.17.112.21",
		"&bridge({origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/external/"+fromNumber+"@10.17.112.21)"))
	response["jobId"] = eslJobid
	return c.JSON(http.StatusOK, response)
}
