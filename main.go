package main

import (
	"github.com/labstack/echo"
	"echoServer/logger"
	"callServer/configs"
	adapters "callServer/adapters/repository"
	requests "callServer/requests/repository"
	auth "callServer/basicAuth/repository"
	esl "github.com/0x19/goesl"
	"github.com/labstack/echo/middleware"
	"flag"
)

var (
	fshost   = flag.String("fshost", "localhost", "Freeswitch hostname. Default: localhost")
	fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
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

	//ESL
	eslClient, err := esl.NewClient(*fshost, *fsport, *password, *timeout)
	if err != nil {
		panic("not able to connect with FreeSWITCH")
		return
	}
	go eslClient.Handle()


	//Associating the controller
	auth.NewAuthController(e)
	requests.NewRequestController(e, eslClient)

	if err := e.Start(config.HttpConfig.HostPort); err != nil {
		log.WithError(err).Fatal("echo server not able to start")
	}
}
