package model

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	sqlite "../external/gosqlite"
	"../utils"
)

var logger = utils.Logger("model")

type Voicemail struct {
	Id            int
	Caller        string
	Called        string
	Date          time.Time
	Duration      time.Duration
	VoicemailPath string
}

type Database struct {
	channel    chan func(db *sqlite.Conn)
	storageDir string
}

func OpenDatabase(dbFile string, storageDir string) Database {
	ch := make(chan func(db *sqlite.Conn))
	go func() {
		db, err := sqlite.Open(dbFile,
			sqlite.OpenReadWrite,
			sqlite.OpenCreate,
			sqlite.OpenFullMutex)

		if err != nil {
			logger.Panic(err)
		}
		defer db.Close()

		db.Exec(`CREATE TABLE IF NOT EXISTS voicemail (
                     id INTEGER PRIMARY KEY,
                     caller TEXT,
                     called TEXT,
                     date INTEGER,
                     duration INTEGER,
                     voicemail TEXT);`)

		for {
			f := <-ch
			f(db)
		}
	}()
	return Database{channel: ch, storageDir: storageDir}
}

func (db Database) GetVoicemails(limit int) ([]Voicemail, error) {
	query := "SELECT * FROM voicemail ORDER BY date DESC LIMIT " + strconv.Itoa(limit)
	errorChannel := make(chan error)

	voicemails := []Voicemail{}

	db.channel <- func(db *sqlite.Conn) {
		s, err := db.Prepare(query)
		if err != nil {
			errorChannel <- err
			return
		}
		defer s.Finalize()

		err = s.Select(func(s *sqlite.Stmt) (err error) {
			var voicemail Voicemail
			var duration string
			if err = s.NamedScan(
				"id", &voicemail.Id,
				"caller", &voicemail.Caller,
				"called", &voicemail.Called,
				"date", &voicemail.Date,
				"duration", &duration,
				"voicemail", &voicemail.VoicemailPath); err != nil {

				return err
			}

			if voicemail.Caller == "" {
				voicemail.Caller = "Unbekannt"
			}
			if voicemail.Called == "" {
				voicemail.Called = "Unbekannt"
			}

			voicemail.Duration, _ = time.ParseDuration(duration + "s")
			voicemails = append(voicemails, voicemail)

			return err
		})

		errorChannel <- err
	}

	return voicemails, <-errorChannel
}

func createFile(dir string, filenameBase, filenameExt string) (string, *os.File, error) {
	var (
		file     *os.File
		err      error
		maxTries = 50
	)

	filename := path.Join(dir, filenameBase+filenameExt)
	for i := 0; i < maxTries; i++ {
		if file, err = os.OpenFile(filename,
			os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600); err == nil {
			break
		}

		filename = path.Join(dir, fmt.Sprintf("%s-%d.mp3", filenameBase, i))
	}

	if file == nil || err != nil {
		return "", nil, err
	}

	return filename, file, nil
}

func saveVoicemailAudio(dir string, voicemail []byte) (string, error) {
	filenameBase := time.Now().Format("20060102-150405")
	filename, file, err := createFile(dir, filenameBase, ".mp3")
	if err == nil {
		file.Write(voicemail)
		file.Close()
	}

	return filename, err
}

func (db Database) AddVoicemail(voicemail Voicemail, voicemailAudio []byte) error {
	errorChannel := make(chan error)

	db.channel <- func(conn *sqlite.Conn) {
		voicemailPath, err := saveVoicemailAudio(db.storageDir, voicemailAudio)
		if err != nil {
			logger.Print("Unable to save voicemail as MP3: ", err)
			errorChannel <- err
			return
		}

		voicemail.VoicemailPath = path.Base(voicemailPath)
		logger.Print("Voicemail saved to ", voicemail.VoicemailPath)

		ins, err := conn.Prepare(`INSERT INTO voicemail (
                                caller, called, date, duration, voicemail
                            ) VALUES (?, ?, ?, ?, ?)`)
		if err != nil {
			errorChannel <- err
			return
		}
		defer ins.Finalize()

		_, err = ins.Insert(voicemail.Caller,
			voicemail.Called,
			voicemail.Date,
			voicemail.Duration.Seconds(),
			voicemail.VoicemailPath)
		if err != nil {
			errorChannel <- err
			return
		}

		errorChannel <- nil
	}

	return <-errorChannel
}
