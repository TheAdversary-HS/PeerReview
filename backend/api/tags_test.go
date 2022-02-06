package api

import (
	"TheAdversary/database"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTags(t *testing.T) {
	if err := initTestDatabase("tags_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	tags := []string{
		"test",
		"linux",
	}

	database.GetDB().Table("article").Create([]map[string]interface{}{
		{
			"title":    "Upload test",
			"summary":  "An example article to test the upload api endpoint",
			"created":  time.Now().Unix(),
			"link":     "article/upload-test",
			"markdown": "Oh god i have to test all this, what am i doing with my life",
			"html":     "<p>Oh god i have to test all this, what am i doing with my life<p>",
		},
	})
	database.GetDB().Table("author").Create([]map[string]interface{}{
		{
			"name":        "me",
			"password":    "",
			"information": "this is my account",
		},
	})
	database.GetDB().Table("article_author").Create([]map[string]interface{}{
		{
			"article_id": 1,
			"author_id":  1,
		},
	})
	database.GetDB().Table("article_tag").Create([]map[string]interface{}{
		{
			"article_id": 1,
			"tag":        "test",
		},
		{
			"article_id": 1,
			"tag":        "linux",
		},
	})

	server := httptest.NewServer(http.HandlerFunc(Tags))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method:     http.MethodGet,
			ResultBody: tags,
			Code:       http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"name": "abc",
			},
			ResultBody: []interface{}{},
			Code:       http.StatusOK,
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
