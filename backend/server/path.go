package server

import (
	"TheAdversary/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ServePath(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(config.FrontendDir, strings.TrimPrefix(r.URL.Path, config.SubPath))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		Error404(w, r)
	} else {
		http.ServeFile(w, r, path)
	}
}
