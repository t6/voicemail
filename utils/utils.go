package utils

import (
	sqlite "github.com/gwenn/gosqlite"
)

func OpenDatabase(dbFile string) (*sqlite.Conn, error) {
	db, err := sqlite.Open(dbFile,
		sqlite.OPEN_READWRITE,
		sqlite.OPEN_CREATE,
		sqlite.OPEN_FULLMUTEX)
	return db, err
}

