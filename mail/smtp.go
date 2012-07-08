package mail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/mail"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	sqlite "github.com/gwenn/gosqlite"
	"github.com/sloonz/go-qprintable"

	. "voicemail/utils"
)

func handleData(rd bufio.Reader, out io.Writer) string {
	// Stolen from go-smtpd:
	// http://code.google.com/p/go-smtpd/source/browse/smtpd/smtpd.go
	fmt.Fprint(out, "354 Go ahead\r\n")

	buf := new(bytes.Buffer)

	for {
		sl, err := rd.ReadSlice('\n')
		if err != nil {
			return ""
		}
		if bytes.Equal(sl, []byte(".\r\n")) {
			break
		}
		if sl[0] == '.' {
			sl = sl[1:]
		}
		_, err = buf.Write(sl)
		if err != nil {
			return ""
		}
	}

	fmt.Fprint(out, "250 2.0.0 Ok: queued\r\n")

	return buf.String()
}

func receiveMessage(in bufio.Reader, out io.Writer) (string, error) {
	// HACK: Here we should negotiate with the SMTP
	// client and tell it what capabilities we have.
	// But we instead skip the whole EHLO and HELO thing
	// and hope for the best!
	fmt.Fprint(out, "220 voicemail ESMTP voicemail\r\n")
	fmt.Fprint(out, "250-voicemail\r\n")
	fmt.Fprint(out, "250 DSN\r\n")

	for {
		line, _, err := in.ReadLine()
		if err != nil {
			return "", err
		}

		cmd := strings.Split(string(line), " ")
		switch cmd[0] {
		case "MAIL", "RCPT":
			fmt.Fprint(out, "250 2.1.0 Ok\r\n")
		case "DATA":
			return handleData(in, out), nil
		case "QUIT":
			return "", errors.New("No message received!")
		}
	}

	return "", errors.New("No message received!")
}

func extractPart(rawMsg string, contentType string, boundary string) ([]byte, error) {
	msg, err := mail.ReadMessage(bytes.NewBufferString(rawMsg))
	if err != nil {
		return nil, err
	}

	parts := multipart.NewReader(msg.Body, boundary)
	part, err := parts.NextPart()
	for err == nil {
		type_ := strings.Split(part.Header["Content-Type"][0], ";")[0]
		if contentType == type_ {
			body, _ := ioutil.ReadAll(part)
			return body, nil
		}
		part, err = parts.NextPart()
	}

	return nil, errors.New("No part found!")
}

func unquote(quotedText []byte) string {
	unquoted, _ := ioutil.ReadAll(qprintable.NewDecoder(
		qprintable.UnixTextEncoding, bytes.NewBuffer(quotedText)))
	return string(unquoted)
}

func extractCall(msg string) (*Call, error) {
	metadata, err := extractPart(msg,
		"text/html",
		"==AVM_Fritz_Box==multipart/alternative==1==")
	if err != nil {
		return nil, err
	}

	html := unquote(metadata)

	split := strings.Split(html, "\"c3\"")

	caller := strings.TrimSpace(strings.Split(split[1][1:], "<")[0])
	called := strings.TrimSpace(strings.Split(split[2][1:], "<")[0])

	duration, err := time.ParseDuration(
		strings.Replace(strings.TrimSpace(split[5][1:6]), ":", "h", 1) + "m")
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2.01.06 15:04 -0700",
		strings.TrimSpace(strings.Split(split[3][1:], "<")[0]+" "+split[4][1:6]) + " +0200")
	if err != nil {
		return nil, err
	}

	return &Call{
		Caller:   caller,
		Called:   called,
		Date:     date,
		Duration: duration,
	}, nil
}

func extractVoicemail(msg string) ([]byte, error) {
	voicemailB64, err := extractPart(msg, "audio/x-wav",
		"==AVM_Fritz_Box==multipart/mixed==0==")
	if err != nil {
		return nil, err
	}
	voicemail := base64.NewDecoder(base64.StdEncoding,
		bytes.NewBuffer(bytes.Replace(
			voicemailB64, []byte("\r\n"), []byte{}, -1)))

	// For the voicemail to be useful, it needs to be a MP3 file.
	// The voicemails are encoded with the aLaw/uLaw codec.
	// lame does not know about it, so we use mplayer to decode
	// the file and then pipe it to lame.
	mplayer := exec.Command("mplayer",
		"-really-quiet",
		"-cache", "8192",
		"-ao", "pcm:waveheader:fast:file=/dev/stdout",
		"-vc", "null",
		"-vo", "null",
		"-",
	)
	lame := exec.Command("lame", "-b16", "-", "-")

	var lameOut bytes.Buffer
	mplayer.Stdin = voicemail
	mplayerOut, err := mplayer.StdoutPipe()
	if err != nil {
		return nil, err
	}
	lame.Stdin = mplayerOut
	lame.Stdout = &lameOut

	if err := mplayer.Start(); err != nil {
		return nil, err
	}
	if err := lame.Run(); err != nil {
		return nil, err
	}

	return lameOut.Bytes(), nil
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

func saveVoicemail(dir string, voicemail []byte) (string, error) {
	filenameBase := time.Now().Format("20060102-150405")
	filename, file, err := createFile(dir, filenameBase, ".mp3")
	if err == nil {
		file.Write(voicemail)
		file.Close()
	}

	return filename, err
}

func createTable(db *sqlite.Conn) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS voicemail (
                        id INTEGER PRIMARY KEY,
                        caller TEXT,
                        called TEXT,
                        date INTEGER,
                        duration INTEGER,
                        voicemail TEXT);`)
}

func insertCall(db *sqlite.Conn, call *Call) error {
	ins, err := db.Prepare(`INSERT INTO voicemail (
                                caller, called, date, duration, voicemail
                            ) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	defer ins.Finalize()

	_, err = ins.Insert(call.Caller,
		call.Called,
		call.Date,
		call.Duration.Seconds(),
		call.VoicemailPath)
	if err != nil {
		return err
	}

	return nil
}

func ProcessMessage(conn net.Conn, storageDir string) (*Call, error) {
	in := bufio.NewReader(conn)
	msg, err := receiveMessage(*in, conn)
	if err != nil {
		return nil, err
	}

	call, err := extractCall(msg)
	if err != nil {
		return nil, err
	}
	log.Print("Received new voicemail ", call)

	voicemail, err := extractVoicemail(msg)
	if err != nil {
		return nil, err
	}

	voicemailPath, err := saveVoicemail(storageDir, voicemail)
	if err != nil {
		log.Print("Unable to save voicemail as MP3: ", err)
		return nil, err
	}

	call.VoicemailPath = path.Base(voicemailPath)
	log.Print("Voicemail saved to ", call.VoicemailPath)

	return call, nil
}

func SaveToDB(dbFile string, call *Call) error {
	db, err := OpenDatabase(dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		return err
	}

	if err := insertCall(db, call); err != nil {
		return err
	}

	return nil
}

func Serve(l net.Listener, dbFile, mp3Dir string) {
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		log.Print("Incoming voicemail from ", conn.RemoteAddr(), "?")
		call, err := ProcessMessage(conn, mp3Dir)
		if err != nil {
			log.Print("Unable to process voicemail: ", err)
		} else {
			if err := SaveToDB(dbFile, call); err != nil {
				log.Print("Unable to save to database: ", err)
			} else {
				log.Print("Save to database successfull.")
			}
		}
		conn.Close()
	}
}
