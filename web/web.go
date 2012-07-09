package web

import (
	"fmt"
	"log"
	"path"
	"math"
	"net"
	"time"
	"net/http"
	"html/template"

	sqlite "github.com/gwenn/gosqlite"

	. "bitbucket.org/tobik/voicemail/utils"

	"bitbucket.org/tobik/voicemail/web/assets"
)

var db *sqlite.Conn
var rootTemplate *template.Template
var voicemailDir string

func isNewMessage(t time.Time) bool {
	// All messages that are 48 hours old are new messages
	return math.Abs(time.Now().Sub(t).Hours()) < 48;
}

type Group struct {
	New []Call
	Old []Call
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	s, err := db.Prepare("SELECT * from voicemail ORDER BY date DESC")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return;
	}
	defer s.Finalize()

	newMessageGroup := []Call{}
	oldMessageGroup := []Call{}
	
	err = s.Select(func(s *sqlite.Stmt) (err error) {
		var call Call
		var duration string
		var voicemailPath string
		if err = s.NamedScan(
			"id", &call.Id,
			"caller", &call.Caller,
			"called", &call.Called,
			"date", &call.Date,
			"duration", &duration,
			"voicemail", &voicemailPath); err != nil {
			return err
		}

		if call.Caller == "" { call.Caller = "Unbekannt" }
		if call.Called == "" { call.Called = "Unbekannt" }
		
		call.VoicemailPath = path.Join("/voicemail", voicemailPath)
		call.Duration, _ = time.ParseDuration(duration + "s")
		if isNewMessage(call.Date) {
			newMessageGroup = append(newMessageGroup, call)
		} else {
			oldMessageGroup = append(oldMessageGroup, call)
		}

		return err
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	err = rootTemplate.ExecuteTemplate(w, "calls", Group{newMessageGroup, oldMessageGroup})
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return;
	}
}

func handleAsset(f func() []byte, t string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", t)
		w.Write(f())
	}
}

func Serve(l net.Listener, dbFile string, voicemailDir string) {
	db, _ = OpenDatabase(dbFile)
	defer db.Close()

	rootTemplate = template.New("root")
	_, err := rootTemplate.Parse(string(app_html()))
	if err != nil {
		panic(err)
	}
	
	http.Handle("/voicemail/", http.StripPrefix("/voicemail/",
		http.FileServer(http.Dir(voicemailDir))))

	http.HandleFunc("/js", handleAsset(assets.Jquery_min_js, "text/javascript"))
	http.HandleFunc("/css", handleAsset(assets.Bootstrap_combined_min_css, "text/css"))
	http.HandleFunc("/img/glyphicons-halflings.png",
		handleAsset(assets.Glyphicons_halflings_png, "image/png"))
	http.HandleFunc("/img/glyphicons-halflings-white.png",
		handleAsset(assets.Glyphicons_halflings_white_png, "image/png"))
	
	http.HandleFunc("/", rootHandler)

	log.Fatal(http.Serve(l, nil))
}
