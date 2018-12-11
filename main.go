package main

import (
	"github.com/labstack/echo"
	"echoServer/logger"
	"callServer/configs"
	adapters "callServer/adapters/repository"
	requests "callServer/requests/repository"
	auth "callServer/basicAuth/repository"
	fs "callServer/eslAdapter/repository"
	utils "callServer/coreUtils/repository"
	"github.com/labstack/echo/middleware"
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

	connId, err := utils.GenUUID()
	if err != nil{
		panic("not able to generate the connection UUID to connect with FreeSWITCH")
	}
	log.Info("uuid got generated for freeswutch" + connId)
	//Associating the controller
	auth.NewAuthController(e)
	requests.NewRequestController(e)

	if err := e.Start(config.HttpConfig.HostPort); err != nil {
		log.WithError(err).Fatal("echo server not able to start")
	}
}
