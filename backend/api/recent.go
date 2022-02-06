package api

import (
	"TheAdversary/config"
	"TheAdversary/database"
	"TheAdversary/schema"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

func Recent(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	request := database.GetDB().Table("article")

	limit := 20
	if query.Has("limit") {
		var err error
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			ApiError{"invalid 'limit' parameter", http.StatusUnprocessableEntity}.Send(w)
			return
		} else if limit > 100 {
			ApiError{"'limit' parameter must not be over 100", http.StatusUnprocessableEntity}.Send(w)
			return
		}
	}
	request.Limit(limit)

	var articleSummaries []schema.ArticleSummary
	request.Find(&articleSummaries)

	for i, summary := range articleSummaries {
		database.GetDB().Table("author").Where("id IN (?)", database.GetDB().Table("article_author").Select("author_id").Where("article_id = ?", summary.Id)).Find(&summary.Authors)
		summary.Tags = []string{}
		database.GetDB().Table("article_tag").Select("tag").Where("article_id = ?", summary.Id).Find(&summary.Tags)
		summary.Link = path.Join(config.SubPath, "article", summary.Link)
		articleSummaries[i] = summary
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articleSummaries)
}
