package repository

import (
	"net/http"
	"github.com/labstack/echo"
	"callServer/adapters"
	"fmt"
	"strings"
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
	fromNumber := callDetails["To"].(string)
	toNumber := callDetails["From"].(string)
	didNumber := callDetails["From"].(string)
	trunkIP := "115.248.91.197"
	response["message"] = "Call is Created"

	originateCommand := fmt.Sprintf("originate %s %s",
		"{origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/internal/"+toNumber+"@"+trunkIP,
		"&bridge({origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/external/"+fromNumber+"@"+trunkIP+")")
	resp,_ := a.ESLClient.SendBgApiCmd(originateCommand)
	respField := strings.Fields(resp)
	fmt.Println(respField)
	response["UUID"] = string(respField[0])
	return c.JSON(http.StatusOK, response)
}
