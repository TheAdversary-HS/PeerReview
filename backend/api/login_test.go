package api

import (
	"TheAdversary/database"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	if err := initTestDatabase("login_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("super secure password"), bcrypt.DefaultCost)
	database.GetDB().Table("author").Create([]map[string]interface{}{
		{
			"name":        "owner",
			"password":    password,
			"information": "owner of the best blog in the world",
		},
	})

	server := httptest.NewServer(http.HandlerFunc(Login))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodPost,
			Body: loginPayload{
				Username: "owner",
				Password: "super secure password",
			},
			ResultCookie: []string{"session_id"},
			Code:         http.StatusOK,
		},
		{
			Method: http.MethodPost,
			Body: loginPayload{
				Username: "not a user",
			},
			Code: http.StatusUnauthorized,
		},
		{
			Method: http.MethodPost,
			Body: loginPayload{
				Username: "test",
			},
			Code: http.StatusUnauthorized,
		},
	})
}
