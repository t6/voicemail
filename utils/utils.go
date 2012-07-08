package utils

import (
	"time"

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
		sqlite.OPEN_READWRITE,
		sqlite.OPEN_CREATE,
		sqlite.OPEN_FULLMUTEX)
	return db, err
}
