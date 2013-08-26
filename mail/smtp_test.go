package mail

import (
	"io/ioutil"
	"net"
	"net/smtp"
	"os"
	"path"
	"testing"
	"time"

	. "bitbucket.org/tobik/voicemail/model"
)

func sendMail(t *testing.T, unixSocket, testData string) {
	conn, err := net.Dial("unix", unixSocket)
	if err != nil {
		t.Error(err)
	}

	c, err := smtp.NewClient(conn, "localhost")
	defer c.Quit()

	if err != nil {
		t.Error(err)
	}

	c.Mail("ignored")
	c.Rcpt("ignored")

	w, err := c.Data()
	defer w.Close()
	if err != nil {
		t.Error(err)
	}

	data, err := ioutil.ReadFile(testData)
	if err != nil {
		t.Error(err)
	}

	w.Write(data)
}

func TestReceive(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "voicemail")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(tempDir)

	unixSocket := path.Join(tempDir, "unix")

	l, err := net.Listen("unix", unixSocket)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer l.Close()

	go sendMail(t, unixSocket, "smtp_test.data")

	conn, err := l.Accept()
	if err != nil {
		t.Error(err)
	}

	voicemail, _, err := ProcessMessage(conn, tempDir)
	if err != nil {
		t.Error(err)
	}

	duration, _ := time.ParseDuration("3s")
	referenceVoicemail := &Voicemail{
		Caller:        "5552341222",
		Called:        "12312234",
		Date:          time.Unix(1256211300-7200, 0),
		Duration:      duration,
		VoicemailPath: voicemail.VoicemailPath,
	}

	_, err = os.Stat(voicemail.VoicemailPath)
	if err == os.ErrNotExist ||
		!voicemail.Date.Equal(referenceVoicemail.Date) ||
		voicemail.Caller != referenceVoicemail.Caller ||
		voicemail.Duration.String() != referenceVoicemail.Duration.String() {
		t.Errorf("Voicemail garbled: expected: %v, actual: %v\n", referenceVoicemail, voicemail)
	}
}
