package api

import (
	"TheAdversary/config"
	"TheAdversary/database"
	"TheAdversary/schema"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gomarkdown/markdown"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func Article(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		articleGet(w, r)
	case http.MethodPost:
		articlePost(w, r)
	case http.MethodDelete:
		articleDelete(w, r)
	case http.MethodPatch:
		articlePatch(w, r)
	}
}

type getResponse struct {
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Authors []int    `json:"authors"`
	Image   string   `json:"image"`
	Tags    []string `json:"tags"`
	Link    string   `json:"link"`
	Content string   `json:"content"`
}

func articleGet(w http.ResponseWriter, r *http.Request) {
	authorId, ok := authorizedSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawId := r.URL.Query().Get("id")

	if rawId == "" {
		ApiError{Message: "no id was given", Code: http.StatusBadRequest}.Send(w)
		return
	}
	id, err := strconv.Atoi(rawId)
	if err != nil {
		ApiError{"invalid 'id' parameter", http.StatusUnprocessableEntity}.Send(w)
		return
	}

	if !database.Exists(database.GetDB().Table("article_author"), "author_id=?", authorId) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var article database.Article
	if database.GetDB().Table("article").First(&article, id).RowsAffected == 0 {
		ApiError{Message: "no such id", Code: http.StatusNotFound}.Send(w)
		return
	}

	resp := getResponse{
		Title:   article.Title,
		Summary: article.Summary,
		Image:   article.Image,
		Link:    article.Link,
		Content: base64.StdEncoding.EncodeToString([]byte(article.Markdown)),
	}
	database.GetDB().Table("article_author").Select("author_id").Where("article_id", article.Id).Find(&resp.Authors)
	database.GetDB().Table("article_tag").Select("tag").Where("article_id", article.Id).Find(&resp.Tags)

	json.NewEncoder(w).Encode(resp)
}

type uploadPayload struct {
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Authors []int    `json:"authors"`
	Image   string   `json:"image"`
	Tags    []string `json:"tags"`
	Link    string   `json:"link"`
	Content string   `json:"content"`
}

func articlePost(w http.ResponseWriter, r *http.Request) {
	authorId, ok := authorizedSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload uploadPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		InvalidJson.Send(w)
		return
	}

	rawMarkdown, err := base64.StdEncoding.DecodeString(payload.Content)
	if err != nil {
		zap.S().Warnf("Cannot decode base64")
		ApiError{Message: "invalid base64 content", Code: http.StatusUnprocessableEntity}.Send(w)
		return
	}

	var notFound []string
	database.GetDB().Select("* FROM ? EXCEPT ?", payload.Authors, database.GetDB().Table("author").Select("id")).Find(&notFound)
	if len(notFound) > 0 {
		ApiError{fmt.Sprintf("no authors with the id(s) %s were found", strings.Join(notFound, ", ")), http.StatusUnprocessableEntity}.Send(w)
		return
	}

	a := article{
		Title:    payload.Title,
		Summary:  payload.Summary,
		Image:    payload.Image,
		Created:  time.Now().Unix(),
		Modified: time.Unix(0, 0).Unix(),
		Link:     payload.Link,
		Markdown: string(rawMarkdown),
		Html:     string(markdown.ToHTML(rawMarkdown, nil, nil)),
	}
	database.GetDB().Table("article").Create(&a)
	var authors []map[string]interface{}
	for _, author := range append([]int{authorId}, payload.Authors...) {
		authors = append(authors, map[string]interface{}{
			"article_id": a.ID,
			"author_id":  author,
		})
	}
	database.GetDB().Table("article_author").Create(&authors)
	if len(payload.Tags) > 0 {
		var tags []map[string]interface{}
		for _, tag := range payload.Tags {
			authors = append(authors, map[string]interface{}{
				"article_id": a.ID,
				"tag":        tag,
			})
		}
		database.GetDB().Table("article_tag").Create(&tags)
	}

	var articleSummary schema.ArticleSummary
	database.GetDB().Table("article").Find(&articleSummary, &a.ID)
	database.GetDB().Table("author").Where("id IN (?)", database.GetDB().Table("article_author").Select("author_id").Where("article_id = ?", a.ID)).Find(&articleSummary.Authors)
	if payload.Tags != nil {
		articleSummary.Tags = payload.Tags
	} else {
		articleSummary.Tags = []string{}
	}
	articleSummary.Link = path.Join(config.SubPath, "article", url.PathEscape(articleSummary.Link))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(articleSummary)
}

type editPayload struct {
	Id      int       `json:"id"`
	Title   *string   `json:"title"`
	Summary *string   `json:"summary"`
	Authors *[]int    `json:"authors"`
	Image   *string   `json:"image"`
	Tags    *[]string `json:"tags"`
	Link    *string   `json:"link"`
	Content *string   `json:"content"`
}

func articlePatch(w http.ResponseWriter, r *http.Request) {
	authorId, ok := authorizedSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload editPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		InvalidJson.Send(w)
		return
	}
	if !database.Exists(database.GetDB().Table("article"), "id = ?", payload.Id) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if !database.Exists(database.GetDB().Table("article_author"), "article_id = ? AND author_id = ?", payload.Id, authorId) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	updates := map[string]interface{}{}
	var authorUpdates, tagUpdates []map[string]interface{}

	a := article{
		ID:       payload.Id,
		Modified: time.Now().Unix(),
	}

	if payload.Title != nil {
		updates["title"] = *payload.Title
	}
	if payload.Summary != nil {
		updates["summary"] = *payload.Summary
	}
	if payload.Authors != nil {
		var notFound []string
		database.GetDB().Select("* FROM ? EXCEPT ?", *payload.Authors, database.GetDB().Table("author").Select("id")).Find(&notFound)
		if len(notFound) > 0 {
			ApiError{fmt.Sprintf("no authors with the id(s) %s were found", strings.Join(notFound, ", ")), http.StatusUnprocessableEntity}.Send(w)
			return
		}
		for _, author := range append([]int{authorId}, *payload.Authors...) {
			authorUpdates = append(authorUpdates, map[string]interface{}{
				"article_id": payload.Id,
				"author_id":  author,
			})
		}
	}
	if payload.Tags != nil {
		for _, tag := range *payload.Tags {
			tagUpdates = append(tagUpdates, map[string]interface{}{
				"article_id": payload.Id,
				"tag":        tag,
			})
		}
	}
	if payload.Image != nil {
		updates["image"] = *payload.Image
	}
	if payload.Link != nil {
		updates["link"] = *payload.Link
	}
	if payload.Content != nil {
		rawMarkdown, err := base64.StdEncoding.DecodeString(*payload.Content)
		if err != nil {
			zap.S().Warnf("Cannot decode base64")
			ApiError{Message: "invalid base64 content", Code: http.StatusUnprocessableEntity}.Send(w)
			return
		}
		a.Markdown = string(rawMarkdown)
		a.Html = string(markdown.ToHTML(rawMarkdown, nil, nil))
		updates["markdown"] = string(rawMarkdown)
		updates["html"] = string(markdown.ToHTML(rawMarkdown, nil, nil))
	}

	if len(updates) > 0 {
		updates["modified"] = time.Now().Unix()

		database.GetDB().Table("article").Where("id = ?", payload.Id).Updates(&updates)
		if authorUpdates != nil {
			database.GetDB().Table("article_author").Where("article_id = ?", payload.Id).Delete(nil)
			database.GetDB().Table("article_author").Create(&authorUpdates)
		}
		if tagUpdates != nil {
			database.GetDB().Table("article_tag").Where("article_id = ?", payload.Id).Delete(nil)
			database.GetDB().Table("article_tag").Create(&tagUpdates)
		}
	}

	var articleSummary schema.ArticleSummary
	database.GetDB().Table("article").Find(&articleSummary, payload.Id)
	database.GetDB().Table("author").Where("id IN (?)", database.GetDB().Table("article_author").Select("author_id").Where("article_id = ?", payload.Id)).Find(&articleSummary.Authors)
	if payload.Tags != nil {
		articleSummary.Tags = *payload.Tags
	} else {
		articleSummary.Tags = []string{}
	}
	articleSummary.Link = path.Join(config.SubPath, "article", url.PathEscape(articleSummary.Link))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articleSummary)
}

type deletePayload struct {
	Id int `json:"id"`
}

func articleDelete(w http.ResponseWriter, r *http.Request) {
	authorId, ok := authorizedSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload deletePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		InvalidJson.Send(w)
		return
	}
	if !database.Exists(database.GetDB().Table("article"), "id = ?", payload.Id) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if !database.Exists(database.GetDB().Table("article_author"), "article_id = ? AND author_id = ?", payload.Id, authorId) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	database.GetDB().Table("article").Delete(&payload)

	w.WriteHeader(http.StatusOK)
}
