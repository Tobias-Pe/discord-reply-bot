package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	Logger = logger.Sugar()
	Logger.Debug("Logger initialised")
}
