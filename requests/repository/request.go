package repository

import (
	"net/http"
	"github.com/labstack/echo"
	coreUtils "callServer/coreUtils/repository"
	"callServer/adapters"
	"fmt"
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

	callUUID, err := coreUtils.GenUUID()
	if err != nil {
		panic("not able to generate the connection UUID to connect with FreeSWITCH")
	}
	response["callUUID"] = callUUID

	//authId := c.Param("auth_id")
	fromNumber := callDetails["fromNumber"].(string)
	toNumber := callDetails["toNumber"].(string)
	didNumber := callDetails["didNumber"].(string)

	response["message"] = "Call is Created"

	originateCommand := fmt.Sprintf("originate %s %s",
		"{origination_uuid="+callUUID+",origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/internal/"+toNumber+"@10.17.112.21",
		"&bridge({origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/external/"+fromNumber+"@10.17.112.21)")
	a.ESLClient.Originate(originateCommand)
	return c.JSON(http.StatusOK, response)
}
