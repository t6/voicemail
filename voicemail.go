package main

import (
	"flag"
	"log"
	"net"
	"os/user"
	"strconv"
	"syscall"

	"bitbucket.org/tobik/voicemail/mail"
	"bitbucket.org/tobik/voicemail/utils"
	"bitbucket.org/tobik/voicemail/web"
)

var logger *log.Logger = utils.Logger("voicemail")

func main() {
	var Hostname, User, DatabaseFile, VoicemailDirectory, HttpPort, SmtpPort string
	var Limit int

	flag.StringVar(&Hostname, "host", "localhost", "Hostname or IP to bind to")
	flag.StringVar(&User, "user", "nobody", "User to drop to after binding")
	flag.StringVar(&DatabaseFile, "database", "./voicemail.sqlite", "Database file location")
	flag.StringVar(&VoicemailDirectory, "voicemail", "./mp3/", "Voicemail storage directory")
	flag.StringVar(&HttpPort, "http-port", "8080", "Port for the HTTP service")
	flag.StringVar(&SmtpPort, "smtp-port", "2500", "Port for the SMTP service")
	flag.IntVar(&Limit, "limit", -1, "Only display this many voicemails in the web interface")

	flag.Parse()

	smtpListener, err := net.Listen("tcp", Hostname+":"+SmtpPort)
	if err != nil {
		logger.Panic(err)
	}

	httpListener, err := net.Listen("tcp", Hostname+":"+HttpPort)
	if err != nil {
		logger.Panic(err)
	}

	u, err := user.Lookup(User)
	if err != nil {
		logger.Panic(err)
	}

	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	if err = syscall.Setgid(gid); err != nil {
		logger.Panic(err)
	}
	if err = syscall.Setuid(uid); err != nil {
		logger.Panic(err)
	}

	db, err := utils.OpenDatabase(DatabaseFile)
	if err != nil {
		logger.Panic(err)
	}

	go mail.Serve(smtpListener, db, VoicemailDirectory)
	web.Serve(httpListener, db, VoicemailDirectory, Limit)
}
