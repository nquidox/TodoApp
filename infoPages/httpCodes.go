package infoPages

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func infoPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	file, err := os.ReadFile("static/infoPage.html")
	if err != nil {
		log.Error(err)
		return
	}

	_, err = w.Write(file)
	if err != nil {
		log.Error(err)
		return
	}
}
