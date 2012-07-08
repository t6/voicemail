package main

import (
	"flag"
	"net"
	"os/user"
	"strconv"
	"syscall"

	"voicemail/mail"
	"voicemail/web"
)

func main() {
	var Hostname, User, DatabaseFile, VoicemailDirectory, HttpPort, SmtpPort string
	flag.StringVar(&Hostname, "host", "localhost", "Hostname or IP to bind to")
	flag.StringVar(&User, "user", "nobody", "User to drop to after binding")
	flag.StringVar(&DatabaseFile, "database", "./voicemail.sqlite", "Database file location")
	flag.StringVar(&VoicemailDirectory, "voicemail", "./mp3/", "Voicemail storage directory")
	flag.StringVar(&HttpPort, "http-port", "8080", "Port for the HTTP service")
	flag.StringVar(&SmtpPort, "smtp-port", "2500", "Port for the SMTP service")

	flag.Parse()

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
