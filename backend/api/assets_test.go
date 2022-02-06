package api

import (
	"TheAdversary/database"
	"TheAdversary/schema"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssetsGet(t *testing.T) {
	if err := initTestDatabase("assets_get_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	database.GetDB().Table("assets").Create([]map[string]interface{}{
		{
			"name": "linux",
			"data": "this should be an image of tux :3",
			"link": "assets/linux",
		},
		{
			"name": "get test",
			"data": "",
			"link": "assets/get-test",
		},
	})

	server := httptest.NewServer(http.HandlerFunc(assetsGet))
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
			Query: map[string]interface{}{
				"q": "linux",
			},
			ResultBody: []schema.Asset{
				{
					Id:   1,
					Name: "linux",
					Link: "assets/linux",
				},
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Query: map[string]interface{}{
				"limit": 1,
			},
			ResultBody: []schema.Asset{
				{
					Id:   1,
					Name: "linux",
					Link: "assets/linux",
				},
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodGet,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			ResultBody: []schema.Asset{
				{
					Id:   1,
					Name: "linux",
					Link: "assets/linux",
				},
				{
					Id:   2,
					Name: "get test",
					Link: "assets/get-test",
				},
			},
			Code: http.StatusOK,
		},
	})
}

func TestAssetsPost(t *testing.T) {
	if err := initTestDatabase("assets_post_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	mw.SetBoundary("test")
	formFile, _ := mw.CreateFormFile("file", "srfwsr")
	formFile.Write([]byte("just a test file"))
	mw.WriteField("name", "test")
	mw.Close()

	fmt.Println(buf.String())

	server := httptest.NewServer(http.HandlerFunc(assetsPost))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodPost,
			Code:   http.StatusUnauthorized,
		},
		{
			Method: http.MethodPost,
			Header: map[string]string{
				"Content-Type": "multipart/form-data; boundary=test",
			},
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Body: buf.Bytes(),
			ResultBody: schema.Asset{
				Id:   1,
				Name: "test",
				Link: "assets/test",
			},
			Code: http.StatusCreated,
		},
		{
			Method: http.MethodPost,
			Header: map[string]string{
				"Content-Type": "multipart/form-data; boundary=test",
			},
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Body: buf.Bytes(),
			Code: http.StatusConflict,
		},
	})
}

func TestAssetsDelete(t *testing.T) {
	if err := initTestDatabase("assets_delete_test.sqlite3"); err != nil {
		t.Fatal(err)
	}

	database.GetDB().Table("assets").Create(map[string]interface{}{
		"name": "example",
		"data": "just a normal string",
		"link": "assets/example",
	})

	server := httptest.NewServer(http.HandlerFunc(assetsDelete))
	checkTestInformation(t, server.URL, []testInformation{
		{
			Method: http.MethodDelete,
			Code:   http.StatusUnauthorized,
		},
		{
			Method: http.MethodDelete,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Body: assetsDeletePayload{
				Id: 1,
			},
			Code: http.StatusOK,
		},
		{
			Method: http.MethodDelete,
			Cookie: map[string]string{
				"session_id": initSession(),
			},
			Body: assetsDeletePayload{
				Id: 69,
			},
			Code: http.StatusNotFound,
		},
	})
}
