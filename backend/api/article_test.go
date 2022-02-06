package api

import (
	"TheAdversary/database"
	"TheAdversary/schema"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestArticleGet(t *testing.T) {
	if err := initTestDatabase("upload_get_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	database.GetDB().Table("article").Create([]map[string]interface{}{
		{
			"title":    "Get test",
			"summary":  "",
			"created":  time.Now().Unix(),
			"link":     "get-test",
			"markdown": "Testing ._.",
			"html":     "<p>Testing ._.<p>",
		},
	})
	database.GetDB().Table("author").Create([]map[string]interface{}{
		{
			"name":        "admin",
			"password":    "",
			"information": "admin",
		},
	})
	database.GetDB().Table("article_author").Create([]map[string]interface{}{
		{
			"article_id": 1,
			"author_id":  1,
		},
	})

	server := httptest.NewServer(http.HandlerFunc(articleGet))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodGet,
			Code:   http.StatusUnauthorized,
		},
		{
			Method: http.MethodGet,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			ResultBody: map[string]interface{}{
				"message": "no id was given",
			},
			Code: http.StatusBadRequest,
		},
		{
			Method: http.MethodGet,
			Query: map[string]interface{}{
				"id": 1,
			},
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			ResultBody: getResponse{
				Title:   "Get test",
				Authors: []int{1},
				Tags:    []string{},
				Link:    "get-test",
				Content: base64.StdEncoding.EncodeToString([]byte("Testing ._.")),
			},
			Code: http.StatusOK,
		},
	})
}

func TestArticlePost(t *testing.T) {
	if err := initTestDatabase("upload_post_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	database.GetDB().Table("article").Create([]map[string]interface{}{
		{
			"title":    "Upload test",
			"summary":  "An example article to test the upload api endpoint",
			"created":  time.Now().Unix(),
			"link":     "upload-test",
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

	server := httptest.NewServer(http.HandlerFunc(articlePost))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodPost,
			Code:   http.StatusUnauthorized,
		},
		{
			Method: http.MethodPost,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Body: uploadPayload{
				Title:   "Testooo",
				Summary: "I have no idea what to put in here",
				Authors: []int{1},
				Link:    "testooo",
				Content: base64.StdEncoding.EncodeToString([]byte("### Testo")),
			},
			ResultBody: schema.ArticleSummary{
				Id:      2,
				Title:   "Testooo",
				Summary: "I have no idea what to put in here",
				Authors: []schema.Author{
					{
						Id:          1,
						Name:        "me",
						Information: "this is my account",
					},
				},
				Tags: []string{},
				Link: "article/testooo",
			},
			Code: http.StatusCreated,
			AfterExec: func(information *testInformation) {
				var created int64
				database.GetDB().Table("article").Select("created").Where("id = 2").Find(&created)
				res := information.ResultBody.(schema.ArticleSummary)
				res.Created = created
				information.ResultBody = res
			},
		},
	})
}

func TestArticlePatch(t *testing.T) {
	if err := initTestDatabase("edit_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	database.GetDB().Table("article").Create([]map[string]interface{}{
		{
			"title":    "test article",
			"summary":  "example summary",
			"created":  time.Now().Unix(),
			"link":     "test-article",
			"markdown": "Just a simple test article",
			"html":     "<p>Just a simple test article<p>",
		},
	})
	database.GetDB().Table("author").Create([]map[string]interface{}{
		{
			"name":     "test",
			"password": "",
		},
		{
			"name":        "admin",
			"password":    "123456",
			"information": "im the admin",
		},
	})
	database.GetDB().Table("article_author").Create([]map[string]interface{}{
		{
			"article_id": 1,
			"author_id":  1,
		},
	})

	server := httptest.NewServer(http.HandlerFunc(articlePatch))
	newTitle := "New title"
	var created int64
	database.GetDB().Table("article").Select("created").Where("id = 1").Find(&created)
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodPost,
			Code:   http.StatusUnauthorized,
		},
		{
			Method: http.MethodPost,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Body: editPayload{
				Id: 69,
			},
			Code: http.StatusNotFound,
		},
		{
			Method: http.MethodPost,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Body: editPayload{
				Id:    1,
				Title: &newTitle,
				Authors: &[]int{
					2,
				},
			},
			ResultBody: schema.ArticleSummary{
				Id:      1,
				Title:   "New title",
				Summary: "example summary",
				Authors: []schema.Author{
					{
						Id:   1,
						Name: "test",
					},
					{
						Id:          2,
						Name:        "admin",
						Information: "im the admin",
					},
				},
				Created: created,
				Tags:    []string{},
				Link:    "article/test-article",
			},
			AfterExec: func(information *testInformation) {
				var modified int64
				database.GetDB().Table("article").Select("modified").Where("id = 1").Find(&modified)
				res := information.ResultBody.(schema.ArticleSummary)
				res.Modified = modified
				information.ResultBody = res
			},
			Code: http.StatusOK,
		},
	})
}

func TestArticleDelete(t *testing.T) {
	if err := initTestDatabase("delete_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	database.GetDB().Table("article").Create([]map[string]interface{}{
		{
			"title":    "test",
			"summary":  "test summary",
			"image":    "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
			"created":  time.Now().Unix(),
			"modified": time.Now().Unix(),
			"link":     "test",
			"markdown": "# Title",
			"html":     "<h1>Title</h1>",
		},
		{
			"title":    "owo",
			"created":  time.Now().Unix(),
			"link":     "owo",
			"markdown": "owo",
			"html":     "<p>owo<p>",
		},
	})
	database.GetDB().Table("author").Create(map[string]interface{}{
		"name":     "test",
		"password": "",
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
	})

	server := httptest.NewServer(http.HandlerFunc(articleDelete))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodPost,
			Code:   http.StatusUnauthorized,
		},
		{
			Method: http.MethodPost,
			Body: deletePayload{
				Id: 1,
			},
			Code: http.StatusUnauthorized,
		},
		{
			Method: http.MethodPost,
			Body: deletePayload{
				Id: 1,
			},
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodPost,
			Body: deletePayload{
				Id: 69,
			},
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Code: http.StatusNotFound,
		},
	})
}
