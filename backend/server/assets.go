package server

import (
	"TheAdversary/database"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"path"
)

func Assets(w http.ResponseWriter, r *http.Request) {
	assetName := mux.Vars(r)["asset"]
	var buf []interface{}
	database.GetDB().Table("assets").Select("data").Find(&buf, "link = ?", assetName)

	if buf == nil {
		Error404(w, r)
	} else {
		data := buf[0].([]byte)
		w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(assetName)))
		w.Write(data)
	}
}
