package mail

import (
	"io/ioutil"
	"net"
	"net/smtp"
	"os"
	"path"
	"testing"
	"time"
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

	duration, _ := time.ParseDuration("180s")
	referenceCall := &Call{
		caller:        "5552341222",
		called:        "12312234",
		date:          time.Unix(1256211300, 0),
		duration:      duration,
		voicemailPath: call.voicemailPath,
	}

	_, err = os.Stat(call.voicemailPath)
	if err == os.ErrNotExist ||
		!call.date.Equal(referenceCall.date) ||
		call.caller != referenceCall.caller ||
		call.duration.String() != referenceCall.duration.String() {
		t.Errorf("Call garbled: expected: %v, actual: %v\n", referenceCall, call)
	}
}
