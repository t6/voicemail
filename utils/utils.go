package utils

import (
	"time"
	"log"
	"log/syslog"
	
	sqlite "github.com/gwenn/gosqlite"
)

type Call struct {
	Id            int
	Caller        string
	Called        string
	Date          time.Time
	Duration      time.Duration
	VoicemailPath string
}

func OpenDatabase(dbFile string) (*sqlite.Conn, error) {
	db, err := sqlite.Open(dbFile,
		sqlite.OpenReadWrite,
		sqlite.OpenCreate,
		sqlite.OpenFullMutex)
	return db, err
}

func Logger(prefix string) *log.Logger {
	logger, err := syslog.NewLogger(syslog.LOG_INFO, log.Lshortfile)
	if err != nil {
		panic(err)
	}

	logger.SetPrefix("[" + prefix + "] ")
	
	return logger;
}
	
