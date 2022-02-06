package main

import (
	"TheAdversary/api"
	"TheAdversary/config"
	"TheAdversary/database"
	"TheAdversary/server"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	var subrouter *mux.Router
	if config.SubPath != "" {
		subrouter = r.PathPrefix(config.SubPath).Subrouter()
	} else {
		subrouter = r
	}

	setupApi(subrouter)
	setupFrontend(subrouter)

	db, err := database.NewSqlite3Connection(config.DatabaseFile)
	if err != nil {
		panic(err)
	}
	database.SetGlobDB(db)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), r); err != nil {
		panic(err)
	}
}

func setupApi(r *mux.Router) {
	r.HandleFunc("/api/login", api.Login).Methods(http.MethodPost)

	r.HandleFunc("/api/authors", api.Authors).Methods(http.MethodGet)
	r.HandleFunc("/api/tags", api.Tags).Methods(http.MethodGet)

	r.HandleFunc("/api/recent", api.Recent).Methods(http.MethodGet)
	r.HandleFunc("/api/search", api.Search).Methods(http.MethodGet)

	r.HandleFunc("/api/article", api.Article).Methods(http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete)

	r.HandleFunc("/api/assets", api.Assets).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)

	r.MethodNotAllowedHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ApiError{Message: "invalid method", Code: http.StatusNotFound}.Send(w)
	}))
	r.NotFoundHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ApiError{Message: "invalid endpoint", Code: http.StatusNotFound}.Send(w)
	}))
}

func setupFrontend(r *mux.Router) {
	r.HandleFunc("/article/{article}", server.Article).Methods(http.MethodGet)
	r.HandleFunc("/assets/{asset}", server.Assets).Methods(http.MethodGet)

	r.PathPrefix("/css/").HandlerFunc(server.ServePath).Methods(http.MethodGet)
	r.PathPrefix("/img/").HandlerFunc(server.ServePath).Methods(http.MethodGet)
	r.PathPrefix("/js/").HandlerFunc(server.ServePath).Methods(http.MethodGet)
	r.PathPrefix("/html/").HandlerFunc(server.ServePath).Methods(http.MethodGet)

	landingpage := template.Must(template.ParseFiles(filepath.Join(config.FrontendDir, "html", "landingpage.gohtml")))
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		landingpage.Execute(w, struct {
			BasePath string
		}{BasePath: config.Address + strings.TrimSuffix(path.Join("/", config.SubPath), "/") + "/"})
	}).Methods(http.MethodGet)

	r.NotFoundHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Error404(w, r)
	}))
}
