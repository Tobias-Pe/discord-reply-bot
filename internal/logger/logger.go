package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	Logger = logger.Sugar()
	Logger.Debug("Logger initialised")
}
