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
	response := make(map[string]interface{})
	response["message"] = "Successfully Authenticated " + auth_id
	eslJobid := a.ESLClient.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))
	response["jobId"] = eslJobid
	return c.JSON(http.StatusOK, response)
}