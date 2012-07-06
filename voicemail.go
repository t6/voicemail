package main

import (
	"os"
	"voicemail/mail"
)

func main() {
	call := mail.ProcessMessage(os.Stdin, os.Stdout, "mp3/")
	mail.SaveToDB("test.sqlite", call)
}
