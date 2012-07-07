package main

import (
	"net"
	"os/user"
	"strconv"
	"syscall"
	"voicemail/mail"
	"voicemail/web"
)

func main() {
	Hostname := "localhost"
	User := "tobias"
	DatabaseFile := "/home/tobias/src/voicemail/test.sqlite"
	VoicemailDirectory := "/home/tobias/src/voicemail/mp3/"
	HttpPort := "8080"
	SmtpPort := "2500"

	smtpListener, err := net.Listen("tcp", Hostname+":"+SmtpPort)
	if err != nil {
		panic(err)
	}

	httpListener, err := net.Listen("tcp", Hostname+":"+HttpPort)
	if err != nil {
		panic(err)
	}

	u, err := user.Lookup(User)
	if err != nil {
		panic(err)
	}

	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	if err = syscall.Setgid(gid); err != nil {
		panic(err)
	}
	if err = syscall.Setuid(uid); err != nil {
		panic(err)
	}

	go mail.Serve(smtpListener, DatabaseFile, VoicemailDirectory)
	web.Serve(httpListener, DatabaseFile, VoicemailDirectory)
}
