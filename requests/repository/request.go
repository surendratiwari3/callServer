package repository

import (
	"net/http"
	"github.com/labstack/echo"
	esl "github.com/0x19/goesl"
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
	auth_id := c.Param("auth_id")
	fromNumber := c.QueryParam("fromNumber")
	toNumber := c.QueryParam("toNumber")
	didNumber := c.QueryParam("didNumber")
	response := make(map[string]interface{})
	response["message"] = "Successfully Authenticated " + auth_id
	eslJobid := a.ESLClient.BgApi(fmt.Sprintf("originate %s %s", "{origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/internal/"+toNumber+"@10.17.112.21", "&bridge({origination_caller_id_number="+didNumber+",absolute_codec_string=PCMU,PCMA}sofia/external/"+fromNumber+"@10.17.112.21)"))
	response["jobId"] = eslJobid
	return c.JSON(http.StatusOK, response)
}