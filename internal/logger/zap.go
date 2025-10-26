package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init() {

	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func Sync() {
	_ = Logger.Sync()
}
