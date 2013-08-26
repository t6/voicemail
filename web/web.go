package web

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net"
	"net/http"
	"path"
	"time"

	. "bitbucket.org/tobik/voicemail/utils"
	"bitbucket.org/tobik/voicemail/model"

	"bitbucket.org/tobik/voicemail/web/assets"
)

var rootTemplate *template.Template
var logger *log.Logger = Logger("web")

func isNewMessage(t time.Time) bool {
	// All messages that are 48 hours old are new messages
	return math.Abs(time.Now().Sub(t).Hours()) < 48
}

type Group struct {
	New []model.Voicemail
	Old []model.Voicemail
}

func rootHandler(db model.Database, limit int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		newMessageGroup := []model.Voicemail{}
		oldMessageGroup := []model.Voicemail{}

		voicemails, err := db.GetVoicemails(limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return			
		}

		for _, voicemail := range voicemails {
			voicemail.VoicemailPath = path.Join("/voicemail", voicemail.VoicemailPath)

			if isNewMessage(voicemail.Date) {
				newMessageGroup = append(newMessageGroup, voicemail)
			} else {
				oldMessageGroup = append(oldMessageGroup, voicemail)
			}
		}

		err = rootTemplate.ExecuteTemplate(w, "calls", Group{newMessageGroup, oldMessageGroup})
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
}

func handleAsset(f func() []byte, t string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", t)
		w.Write(f())
	}
}

func Serve(l net.Listener, db model.Database, voicemailDir string, limit int) {
	rootTemplate = template.New("root")
	_, err := rootTemplate.Parse(string(app_html()))
	if err != nil {
		panic(err)
	}

	http.Handle("/voicemail/", http.StripPrefix("/voicemail/",
		http.FileServer(http.Dir(voicemailDir))))

	http.HandleFunc("/js/zepto.min.js",
		handleAsset(assets.Zepto_min_js, "text/javascript"))
	http.HandleFunc("/css", handleAsset(assets.Bootstrap_combined_min_css, "text/css"))
	http.HandleFunc("/img/glyphicons-halflings.png",
		handleAsset(assets.Glyphicons_halflings_png, "image/png"))
	http.HandleFunc("/img/glyphicons-halflings-white.png",
		handleAsset(assets.Glyphicons_halflings_white_png, "image/png"))
	http.HandleFunc("/img/apple-touch-icon.png",
		handleAsset(assets.Apple_touch_icon_png, "image/png"))

	http.HandleFunc("/", rootHandler(db, limit))

	logger.Fatal(http.Serve(l, nil))
}
