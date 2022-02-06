package api

import (
	"TheAdversary/database"
	"encoding/json"
	"net/http"
	"strconv"
)

func Tags(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	request := database.GetDB().Table("article_tag")

	if query.Has("name") {
		request.Where("tag LIKE ?", "%"+query.Get("name")+"%")
	}
	if query.Has("limit") {
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			ApiError{"invalid 'limit' parameter", http.StatusUnprocessableEntity}.Send(w)
			return
		}
		request.Limit(limit)
	}

	tags := make([]string, 0)
	request.Distinct("tag").Find(&tags)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tags)
}
