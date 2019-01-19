package main

import (
	"fmt"
	"github.com/labstack/echo"
	"callServer/logger"
	"callServer/configs"
	adapters "callServer/adapters/repository"
	requests "callServer/requests/repository"
	auth "callServer/basicAuth/repository"
	"github.com/labstack/echo/middleware"
    "callServer/inboundCallHandler/repository"
)

func main() {

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Setting up the config
	config := configs.GetConfig()

	//Setting up the Logger
	log := logger.NewLogger(config.Log.LogFile, config.Log.LogFile)

	//Setting up the Adapters to use
	redisAdapter, err := adapters.NewRedisAdapterRepository(config)
	if redisAdapter == nil || err != nil {
		log.WithError(err).Fatal("redis is not able to connect ")
		panic("redis is not able to connect")
	}
	eslSessions := adapters.NewESLsessions(config)
	eslAdapter, err := adapters.NewESLAdapterRepository(config,eslSessions)
	if eslAdapter == nil || err != nil {
		log.WithError(err).Fatal("FreeSWITCH is not able to connect ")
		panic("FreeSWITCH is not able to connect")
	}
	fmt.Println("hello world")
	eslEventSessions := repository.NewESLsessions(config)
	go func(){
		repository.NewInboundESLRepository(config,eslEventSessions)
        }()
	fmt.Println("hello world")
	//Associating the controller
	auth.NewAuthController(e)
	requests.NewRequestController(e, eslAdapter)

	if err := e.Start(config.HttpConfig.HostPort); err != nil {
		log.WithError(err).Fatal("echo server not able to start")
	}
}
