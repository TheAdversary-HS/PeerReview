package api

import (
	"TheAdversary/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthors(t *testing.T) {
	if err := initTestDatabase("authors_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	database.GetDB().Table("author").Create([]map[string]interface{}{
		{
			"name":        "test",
			"password":    "",
			"information": "test information",
		},
	})

	server := httptest.NewServer(http.HandlerFunc(Authors))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodGet,
			ResultBody: []map[string]interface{}{
				{
					"id":          1,
					"name":        "test",
					"information": "test information",
				},
			},
			Code: http.StatusOK,
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
