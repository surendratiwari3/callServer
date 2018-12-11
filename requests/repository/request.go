package repository

import (
	"net/http"
	"github.com/labstack/echo"
)

func NewRequestController(e *echo.Echo) {
	e.POST("v1/Account/:auth_id/Call/", requestHandler)
}

func requestHandler(c echo.Context) error {
	auth_id := c.Param("auth_id")
	response := make(map[string]interface{})
	response["message"] = "Successfully Authenticated " + auth_id
	return c.JSON(http.StatusOK, response)
}