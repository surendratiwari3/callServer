package repository

import (
"github.com/labstack/echo"
"github.com/labstack/echo/middleware"
)

func NewAuthController(e *echo.Echo) {
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "joe" && password == "secret" {
			return true, nil
		}
		return false, nil
		// check username and password for a match here
	}))
}