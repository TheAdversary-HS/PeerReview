package api

import (
	"TheAdversary/config"
	"TheAdversary/database"
	"TheAdversary/schema"
	"encoding/json"
	"gorm.io/gorm/clause"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

var assetsPayload struct {
	ArticleId int    `json:"article_id"`
	Content   string `json:"content"`
}

func Assets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		assetsGet(w, r)
	case http.MethodPost:
		assetsPost(w, r)
	case http.MethodDelete:
		assetsDelete(w, r)
	}
}

func assetsGet(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizedSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	request := database.GetDB().Table("assets")

	if query.Has("q") {
		request.Where("LOWER(name) LIKE ?", "%"+query.Get("q")+"%")
	}
	limit := 20
	if query.Has("limit") {
		var err error
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			ApiError{"invalid 'limit' parameter", http.StatusUnprocessableEntity}.Send(w)
			return
		}
	}
	request.Limit(limit)

	var assets []schema.Asset
	request.Find(&assets)

	for _, asset := range assets {
		asset.Link = path.Join(config.SubPath, "assets", asset.Link)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&assets)
}

func assetsPost(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizedSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			ApiError{Message: "file is missing", Code: http.StatusUnprocessableEntity}.Send(w)
		} else {
			ApiError{Message: "could not parse file" + err.Error(), Code: http.StatusInternalServerError}.Send(w)
		}
		return
	}
	defer file.Close()

	var name string
	if name = r.FormValue("name"); name == "" {
		name = header.Filename
	}

	rawData, err := io.ReadAll(file)
	if err != nil {
		ApiError{Message: "failed to read file", Code: http.StatusInternalServerError}.Send(w)
		return
	}

	tmpDatabaseSchema := struct {
		Id   int
		Name string
		Data []byte
		Link string
	}{Name: name, Data: rawData, Link: url.PathEscape(name)}

	if database.GetDB().Table("assets").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&tmpDatabaseSchema).RowsAffected == 0 {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(schema.Asset{
			Id:   tmpDatabaseSchema.Id,
			Name: tmpDatabaseSchema.Name,
			Link: path.Join(config.SubPath, "assets", tmpDatabaseSchema.Link),
		})
	}
}

type assetsDeletePayload struct {
	Id int `json:"id"`
}

func assetsDelete(w http.ResponseWriter, r *http.Request) {
	_, ok := authorizedSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload assetsDeletePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		InvalidJson.Send(w)
		return
	}

	if database.GetDB().Table("assets").Delete(schema.Asset{}, payload.Id).RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
