package mail

import (
	"io/ioutil"
	"net"
	"net/smtp"
	"os"
	"path"
	"testing"
	"time"

	. "bitbucket.org/tobik/voicemail/utils"
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

	call, err := ProcessMessage(conn, tempDir)
	if err != nil {
		t.Error(err)
	}

	duration, _ := time.ParseDuration("3s")
	referenceCall := &Call{
		Caller:        "5552341222",
		Called:        "12312234",
		Date:          time.Unix(1256211300 - 7200, 0),
		Duration:      duration,
		VoicemailPath: call.VoicemailPath,
	}

	_, err = os.Stat(call.VoicemailPath)
	if err == os.ErrNotExist ||
		!call.Date.Equal(referenceCall.Date) ||
		call.Caller != referenceCall.Caller ||
		call.Duration.String() != referenceCall.Duration.String() {
		t.Errorf("Call garbled: expected: %v, actual: %v\n", referenceCall, call)
	}
}
