package api

import (
	"TheAdversary/config"
	"TheAdversary/database"
	"TheAdversary/schema"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	request := database.GetDB().Table("article")

	if query.Has("q") {
		request.Where("LOWER(title) LIKE ?", "%"+query.Get("q")+"%")
	}
	if query.Has("from") {
		from, err := strconv.ParseInt(query.Get("from"), 10, 64)
		if err != nil {
			ApiError{"invalid 'from' parameter", http.StatusUnprocessableEntity}.Send(w)
			return
		}
		request.Where("created >= ?", from).Or("modified >= ?", from)
	}
	if query.Has("to") {
		to, err := strconv.ParseInt(query.Get("to"), 10, 64)
		if err != nil {
			ApiError{"invalid 'to' parameter", http.StatusUnprocessableEntity}.Send(w)
			return
		}
		request.Where("created <= ?", to).Or("modified <= ?", to)
	}
	if query.Has("authors") {
		var authorIds []int
		if err := json.NewDecoder(strings.NewReader(query.Get("authors"))).Decode(&authorIds); err != nil {
			ApiError{"could not parse 'authors' parameter as array of integers / numbers", http.StatusUnprocessableEntity}.Send(w)
			return
		}
		request.Where("id IN (?)", database.GetDB().Table("article_author").Select("article_id").Where("author_id IN (?)", authorIds))
	}
	if query.Has("tags") {
		var tags []string
		if err := json.NewDecoder(strings.NewReader(query.Get("tags"))).Decode(&tags); err != nil {
			ApiError{"could not parse 'tags' parameter as array of strings", http.StatusUnprocessableEntity}.Send(w)
			return
		}
		authorRequest := database.GetDB().Table("article_tag").Select("article_id").Where("tag IN ?", tags)
		request.Where("id IN (?)", authorRequest)
	}
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
