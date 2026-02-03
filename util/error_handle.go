package util

import (
	"ai-test/logger"
	"ai-test/util/level"
)

var log = logger.NewLogger()

func HandleError(message string, err error, errorLevel level.ErrorLevel) {
	if err != nil {
		switch errorLevel {
		case level.INFO:
			log.Infof(message, err)
			return
		case level.WARN:
			log.Warnf(message, err)
			return
		case level.ERROR:
			log.Errorf(message, err)
			return
		case level.FATAL:
			log.Fatalf(message, err)
		}
	}
}
