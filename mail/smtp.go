package mail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/mail"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	sqlite "github.com/gwenn/gosqlite"
	"github.com/sloonz/go-qprintable"
)

type Call struct {
	caller        string
	called        string
	date          time.Time
	duration      time.Duration
	voicemailPath string
}

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

func receiveMessage(in bufio.Reader, out io.Writer) string {

	// HACK: Here we should negotiate with the SMTP
	// client and tell it what capabilities we have.
	// But we instead skip the whole EHLO and HELO thing
	// and hope for the best!
	fmt.Fprint(out, "220 voicemail ESMTP voicemail\r\n")
	fmt.Fprint(out, "250-voicemail\r\n")
	fmt.Fprint(out, "250 DSN\r\n")

	for {
		line, _, _ := in.ReadLine()

		cmd := strings.Split(string(line), " ")
		switch cmd[0] {
		case "MAIL", "RCPT":
			fmt.Fprint(out, "250 2.1.0 Ok\r\n")
		case "DATA":
			return handleData(in, out)
		case "QUIT":
			return ""
		}
	}

	return ""
}

func extractPart(rawMsg string, contentType string, boundary string) []byte {
	msg, err := mail.ReadMessage(bytes.NewBufferString(rawMsg))
	if err != nil {
		panic(err)
	}

	parts := multipart.NewReader(msg.Body, boundary)
	part, err := parts.NextPart()
	for err == nil {
		type_ := strings.Split(part.Header["Content-Type"][0], ";")[0]
		if contentType == type_ {
			body, _ := ioutil.ReadAll(part)
			return body
		}
		part, err = parts.NextPart()
	}

	return nil
}

func unquote(quotedText []byte) string {
	unquoted, _ := ioutil.ReadAll(qprintable.NewDecoder(
		qprintable.UnixTextEncoding, bytes.NewBuffer(quotedText)))
	return string(unquoted)
}

func extractCall(msg string) *Call {
	metadata := extractPart(msg,
		"text/html",
		"==AVM_Fritz_Box==multipart/alternative==1==")
	if metadata == nil {
		os.Exit(1)
	}

	html := unquote(metadata)

	split := strings.Split(html, "\"c3\"")

	caller := strings.TrimSpace(strings.Split(split[1][1:], "<")[0])
	called := strings.TrimSpace(strings.Split(split[2][1:], "<")[0])

	duration, err := time.ParseDuration(
		strings.Replace(strings.TrimSpace(split[5][1:6]), ":", "h", 1) + "m")
	if err != nil {
		panic(err)
	}

	date, err := time.Parse("02.01.06 15:04",
		strings.TrimSpace(split[3][1:9]+" "+split[4][1:6]))
	if err != nil {
		panic(err)
	}

	return &Call{
		caller:   caller,
		called:   called,
		date:     date,
		duration: duration,
	}
}

func extractVoicemail(msg string) []byte {
	voicemailB64 := extractPart(msg, "application/octet-stream",
		"==AVM_Fritz_Box==multipart/mixed==0==")
	if voicemailB64 == nil {
		os.Exit(1)
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
	mplayerOut, _ := mplayer.StdoutPipe()
	lame.Stdin = mplayerOut
	lame.Stdout = &lameOut

	mplayer.Start()
	lame.Run()

	return lameOut.Bytes()
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

func saveVoicemail(dir string, voicemail []byte) string {
	filenameBase := time.Now().Format("20060102-150405")
	filename, file, err := createFile(dir, filenameBase, ".mp3")
	if err == nil {
		file.Write(voicemail)
		file.Close()
	} else {
		panic(err)
	}

	return filename
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

	_, err = ins.Insert(call.caller,
		call.called,
		call.date,
		call.duration.Seconds(),
		call.voicemailPath)
	if err != nil {
		return err
	}

	return nil
}

func ProcessMessage(input io.Reader, output io.Writer, storageDir string) *Call {
	in := bufio.NewReader(input)
	msg := receiveMessage(*in, output)
	call := extractCall(msg)
	voicemail := extractVoicemail(msg)
	call.voicemailPath = saveVoicemail(storageDir, voicemail)

	return call
}

func SaveToDB(dbFile string, call *Call) {
	//fmt.Fprintf(os.Stderr, "%v\n", call)

	db, err := sqlite.Open(dbFile,
		sqlite.OPEN_READWRITE,
		sqlite.OPEN_CREATE,
		sqlite.OPEN_FULLMUTEX)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err := createTable(db); err != nil {
		panic(err)
	}

	if err := insertCall(db, call); err != nil {
		panic(err)
	}
}

