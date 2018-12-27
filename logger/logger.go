package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func NewLogger(fileName string,logLevel string) *logrus.Logger{
	if Logger != nil {
		return Logger
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Logger.Fatal(err)
	}
	Logger = &logrus.Logger{
		Out:       file,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	return Logger
}
