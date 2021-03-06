package utils

import (
	"log"
	"log/syslog"
)

func Logger(prefix string) *log.Logger {
	logger, err := syslog.NewLogger(syslog.LOG_NOTICE, log.Lshortfile)
	if err != nil {
		panic(err)
	}

	logger.SetPrefix("[" + prefix + "] ")

	return logger
}
