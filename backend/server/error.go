package server

import (
	"TheAdversary/config"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func init() {
	// disable default log output because http.ServeFile prints
	// a message if a header is written 2 times or more
	log.Default().SetOutput(io.Discard)
}

func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, filepath.Join(config.FrontendDir, "error", "404.html"))
}

func Error500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	http.ServeFile(w, r, filepath.Join(config.FrontendDir, "error", "500.html"))
}
