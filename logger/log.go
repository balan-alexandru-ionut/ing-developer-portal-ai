package logger

import "go.uber.org/zap"

func NewLogger() *zap.SugaredLogger {
	rawLogger, _ := zap.NewProduction()
	logger := rawLogger.Sugar()

	return logger
}
