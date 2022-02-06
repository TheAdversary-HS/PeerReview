package api

import (
	"TheAdversary/database"
	"TheAdversary/schema"
	"encoding/json"
	"net/http"
	"strconv"
)

func Authors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	request := database.GetDB().Table("author")

	if query.Has("name") {
		request.Where("name LIKE ?", "%"+query.Get("name")+"%")
	}
	if query.Has("limit") {
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			ApiError{"invalid 'limit' parameter", http.StatusUnprocessableEntity}.Send(w)
			return
		}
		request.Limit(limit)
	}

	var authors []schema.Author
	request.Find(&authors)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authors)
}
