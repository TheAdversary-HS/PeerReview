package api

import (
	"TheAdversary/database"
	"TheAdversary/schema"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"
)

func TestRecent(t *testing.T) {
	if err := initTestDatabase("recent_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	authors := []schema.Author{
		{
			Id:       1,
			Name:     "user",
			Password: "",
		},
	}
	articles := []schema.ArticleSummary{
		{
			Id:       1,
			Title:    "test",
			Summary:  "test summary",
			Authors:  authors,
			Tags:     []string{},
			Created:  time.Unix(0, 0).Unix(),
			Modified: time.Now().Unix(),
			Link:     "article/test",
		},
		{
			Id:      2,
			Title:   "Recent Article",
			Summary: "This article is recent",
			Authors: authors,
			Tags:    []string{},
			Created: time.Now().Unix(),
			Link:    "article/recent",
		},
	}

	database.GetDB().Table("article").Create([]map[string]interface{}{
		{
			"title":    articles[0].Title,
			"summary":  articles[0].Summary,
			"created":  articles[0].Created,
			"modified": articles[0].Modified,
			"link":     path.Base(articles[0].Link),
			"markdown": "# Title",
			"html":     "<h1>Title</h1>",
		},
		{
			"title":    articles[1].Title,
			"summary":  articles[1].Summary,
			"created":  articles[1].Created,
			"link":     path.Base(articles[1].Link),
			"markdown": "This is the most recent article",
			"html":     "<p>This is the most recent article</p>",
		},
	})
	database.GetDB().Table("author").Create(authors)
	database.GetDB().Table("article_author").Create([]map[string]interface{}{
		{
			"article_id": 1,
			"author_id":  1,
		},
		{
			"article_id": 2,
			"author_id":  1,
		},
	})

	server := httptest.NewServer(http.HandlerFunc(Recent))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method:     http.MethodGet,
			ResultBody: articles,
			Code:       http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"limit": 10,
			},
			ResultBody: articles,
			Code:       http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"limit": 1001,
			},
			ResultBody: map[string]interface{}{
				"message": "'limit' parameter must not be over 100",
			},
			Code: http.StatusUnprocessableEntity,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"limit": "notanumber",
			},
			ResultBody: map[string]interface{}{
				"message": "invalid 'limit' parameter",
			},
			Code: http.StatusUnprocessableEntity,
		},
	})
}
