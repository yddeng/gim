package logger

import (
	"github.com/yddeng/utils/log"
)

var logger *log.Logger

func InitLogger(basePath string, fileName string) {
	logger = log.NewLogger(basePath, fileName, 1024*1024*2)
	logger.Debugf("%s logger init", fileName)
}

func Logger() *log.Logger {
	return logger
}
