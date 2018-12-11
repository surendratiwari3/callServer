package repository

import (
	"net/http"
	"github.com/labstack/echo"
)

func NewRequestController(e *echo.Echo) {
	e.POST("/auth/*", requestHandler)
}

func requestHandler(c echo.Context) error {
	response := make(map[string]interface{})
	response["message"] = "Successfully Authenticated"
	return c.JSON(http.StatusOK, response)
}