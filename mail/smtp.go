package mail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	log "log"
	"mime/multipart"
	"net"
	"net/mail"
	"os/exec"
	"strings"
	"time"

	"github.com/sloonz/go-qprintable"

	"bitbucket.org/tobik/voicemail/model"
	. "bitbucket.org/tobik/voicemail/utils"
)

var logger *log.Logger = Logger("mail")

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

func extractCall(msg string) (model.Voicemail, error) {
	metadata, err := extractPart(msg,
		"text/html",
		"==AVM_Fritz_Box==multipart/alternative==1==")
	if err != nil {
		return model.Voicemail{}, err
	}

	html := unquote(metadata)

	split := strings.Split(html, "\"c3\"")

	caller := strings.TrimSpace(strings.Split(split[1][1:], "<")[0])
	called := strings.TrimSpace(strings.Split(split[2][1:], "<")[0])

	duration, err := time.ParseDuration(
		strings.Replace(strings.TrimSpace(split[5][1:6]), ":", "m", 1) + "s")
	if err != nil {
		return model.Voicemail{}, err
	}

	date, err := time.Parse("2.01.2006 15:04 -0700",
		strings.TrimSpace(strings.Split(split[3][1:], "<")[0]+" "+split[4][1:6])+" +0200")
	if err != nil {
		return model.Voicemail{}, err
	}

	return model.Voicemail{
		Caller:   caller,
		Called:   called,
		Date:     date,
		Duration: duration,
	}, nil
}

func extractVoicemailAudio(msg string) ([]byte, error) {
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
		mplayer.Process.Kill()
		return nil, err
	}

	return lameOut.Bytes(), nil
}

func ProcessMessage(conn net.Conn) (model.Voicemail, []byte, error) {
	in := bufio.NewReader(conn)
	msg, err := receiveMessage(*in, conn)
	if err != nil {
		return model.Voicemail{}, nil, err
	}

	voicemail, err := extractCall(msg)
	if err != nil {
		return model.Voicemail{}, nil, err
	}
	logger.Print("Received new voicemail ", voicemail)

	voicemailAudio, err := extractVoicemailAudio(msg)
	if err != nil {
		return model.Voicemail{}, nil, err
	}

	return voicemail, voicemailAudio, nil
}

func Serve(l net.Listener, db model.Database) {
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Fatal(err)
			continue
		}

		logger.Print("Incoming voicemail from ", conn.RemoteAddr(), "?")
		voicemail, voicemailAudio, err := ProcessMessage(conn)
		if err != nil {
			logger.Print("Unable to process voicemail: ", err)
		} else {
			if err := db.AddVoicemail(voicemail, voicemailAudio); err != nil {
				logger.Print("Unable to save to database: ", err)
			} else {
				logger.Print("Save to database successful.")
			}
		}
		conn.Close()
	}
}
