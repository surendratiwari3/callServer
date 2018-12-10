package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

func NewLogger(fileName string,logLevel string) *logrus.Logger{
	if logger != nil {
		return logger
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logger.Fatal(err)
	}
	logger = &logrus.Logger{
		Out:       file,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	return logger
}
