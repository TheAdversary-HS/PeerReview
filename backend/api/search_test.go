package api

import (
	"TheAdversary/database"
	"TheAdversary/schema"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	if err := initTestDatabase("search_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	now := time.Now()

	database.GetDB().Table("article").Create([]map[string]interface{}{
		{
			"title":    "First article",
			"created":  now.Unix(),
			"link":     "first-article",
			"markdown": "This is my first article",
			"html":     "<p>This is my first article</p>",
		},
		{
			"title":    "test",
			"summary":  "test summary",
			"image":    "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
			"created":  now.Unix(),
			"modified": now.Add(24 * time.Hour).Unix(),
			"link":     "test",
			"markdown": "# Title",
			"html":     "<h1>Title</h1>",
		},
		{
			"title":    "owo",
			"created":  now.Unix(),
			"modified": now.Add(12 * time.Hour).Unix(),
			"link":     "owo",
			"markdown": "owo",
			"html":     "<p>owo<p>",
		},
	})
	database.GetDB().Table("author").Create([]map[string]interface{}{
		{
			"name":     "test",
			"password": "",
		},
		{
			"name":     "hacr",
			"password": "1234567890",
		},
	})
	database.GetDB().Table("article_tag").Create([]map[string]interface{}{
		{
			"article_id": 1,
			"tag":        "example",
		},
	})
	database.GetDB().Table("article_author").Create([]map[string]interface{}{
		{
			"article_id": 1,
			"author_id":  1,
		},
		{
			"article_id": 2,
			"author_id":  1,
		},
		{
			"article_id": 3,
			"author_id":  2,
		},
	})

	server := httptest.NewServer(http.HandlerFunc(Search))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"q": "first",
			},
			ResultBody: []schema.ArticleSummary{
				{
					Id:      1,
					Title:   "First article",
					Created: now.Unix(),
					Authors: []schema.Author{
						{
							Id:   1,
							Name: "test",
						},
					},
					Tags: []string{
						"example",
					},
					Link: "article/first-article",
				},
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"from": now.Add(1 * time.Hour).Unix(),
				"to":   now.Add(10 * time.Hour).Unix(),
			},
			ResultBody: []schema.ArticleSummary{
				{
					Id:      2,
					Title:   "test",
					Summary: "test summary",
					Authors: []schema.Author{
						{
							Id:   1,
							Name: "test",
						},
					},
					Image:    "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
					Created:  now.Unix(),
					Modified: now.Add(24 * time.Hour).Unix(),
					Tags:     []string{},
					Link:     "article/test",
				},
				{
					Id:    3,
					Title: "owo",
					Authors: []schema.Author{
						{
							Id:   2,
							Name: "hacr",
						},
					},
					Created:  now.Unix(),
					Modified: now.Add(12 * time.Hour).Unix(),
					Tags:     []string{},
					Link:     "article/owo",
				},
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"authors": []int{2},
			},
			ResultBody: []schema.ArticleSummary{
				{
					Id:    3,
					Title: "owo",
					Authors: []schema.Author{
						{
							Id:   2,
							Name: "hacr",
						},
					},
					Created:  now.Unix(),
					Modified: now.Add(12 * time.Hour).Unix(),
					Tags:     []string{},
					Link:     "article/owo",
				},
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"tags": []string{"\"example\""},
			},
			ResultBody: []schema.ArticleSummary{
				{
					Id:    1,
					Title: "First article",
					Authors: []schema.Author{
						{
							Id:   1,
							Name: "test",
						},
					},
					Created: now.Unix(),
					Tags: []string{
						"example",
					},
					Link: "article/first-article",
				},
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"limit": 2,
			},
			ResultBody: []schema.ArticleSummary{
				{
					Id:    1,
					Title: "First article",
					Authors: []schema.Author{
						{
							Id:   1,
							Name: "test",
						},
					},
					Created: now.Unix(),
					Tags: []string{
						"example",
					},
					Link: "article/first-article",
				},
				{
					Id:      2,
					Title:   "test",
					Summary: "test summary",
					Authors: []schema.Author{
						{
							Id:   1,
							Name: "test",
						},
					},
					Image:    "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
					Created:  now.Unix(),
					Modified: now.Add(24 * time.Hour).Unix(),
					Tags:     []string{},
					Link:     "article/test",
				},
			},
			Code: http.StatusOK,
		},
	})
}
