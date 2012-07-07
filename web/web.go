package web

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"voicemail/utils"
)

func Serve(l net.Listener, dbFile string, voicemailDir string) {
	db, _ := utils.OpenDatabase(dbFile)
	defer db.Close()

	http.Handle("/voicemail/",
		http.StripPrefix("/voicemail/",
		http.FileServer(http.Dir(voicemailDir))))

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})

	log.Fatal(http.Serve(l, nil))
}
