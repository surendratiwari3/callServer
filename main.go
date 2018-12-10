package main

import (
	"github.com/labstack/echo"
	"echoServer/logger"
	"callServer/configs"
	adapters "callServer/adapters/repository"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	config := configs.GetConfig()

	log := logger.NewLogger("info.log", "info")

	redisAdapter, err := adapters.NewRedisAdapterRepository(config)

	if redisAdapter==nil || err != nil {
		log.WithError(err).Fatal("redis is not able to connect ")
		panic("redis is not able to connect" )
	}
	if err := e.Start("0.0.0.0:10000"); err != nil {
		log.WithError(err).Fatal("echo server not able to start")
	}
}
