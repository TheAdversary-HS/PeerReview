package api

import (
	"TheAdversary/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type testInformation struct {
	Method string
	Header map[string]string
	Cookie map[string]string
	Body   interface{}
	Query  map[string]interface{}

	ResultBody   interface{}
	ResultCookie []string
	Code         int

	AfterExec func(*testInformation)
}

func initTestDatabase(name string) error {
	path := filepath.Join(os.TempDir(), name)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err = os.Remove(path); err != nil {
			return err
		}
	}

	db, err := database.NewSqlite3Connection(path)
	if err != nil {
		return err
	}
	database.SetGlobDB(db)
	declaration, _ := os.ReadFile("../database.sql")
	db.Exec(string(declaration))

	return nil
}

func initSession() string {
	sessid := sessionId()
	sessions[sessid] = 1
	return sessid
}

func checkTestInformation(t *testing.T, url string, information []testInformation) {
	for i, information := range information {
		var body io.Reader
		if information.Body != nil {
			if b, ok := information.Body.([]byte); ok {
				body = bytes.NewReader(b)
			} else {
				buf, _ := json.Marshal(information.Body)
				body = bytes.NewReader(buf)
			}
		}

		query := url2.Values{}
		if information.Query != nil {
			for key, value := range information.Query {
				query.Add(key, fmt.Sprintf("%v", value))
			}
		}

		req, _ := http.NewRequest(information.Method, fmt.Sprintf("%s?%s", url, query.Encode()), body)
		if information.Cookie != nil {
			for name, value := range information.Cookie {
				req.AddCookie(&http.Cookie{
					Name:  name,
					Value: value,
				})
			}
		}
		if information.Header != nil {
			for name, value := range information.Header {
				req.Header.Set(name, value)
			}
		}

		resp, _ := http.DefaultClient.Do(req)

		if information.AfterExec != nil {
			information.AfterExec(&information)
		}

		if resp.StatusCode != information.Code {
			t.Errorf("Test %d sent invalid status code: expected %d, got %d", i+1, information.Code, resp.StatusCode)
		}

		if resp.Body != nil {
			var respBody, informationBody interface{}
			json.NewDecoder(resp.Body).Decode(&respBody)
			resp.Body.Close()

			tmpInformationBytes, _ := json.Marshal(information.ResultBody)
			json.Unmarshal(tmpInformationBytes, &informationBody)

			if !reflect.DeepEqual(respBody, informationBody) {
				respBytes, _ := json.Marshal(respBody)
				informationBytes, _ := json.Marshal(informationBody)

				// for some reason the maps are sometimes not matched as equal.
				// this is an additional checks if the map bytes are equal
				if !bytes.Equal(respBytes, informationBytes) {
					t.Errorf("Test %d sent invalid response body: expected %s, got %s", i+1, informationBytes, respBytes)
				}
			}
		}

		if information.ResultCookie != nil {
			for _, cookie := range information.ResultCookie {
				var found bool
				for _, respCookie := range resp.Cookies() {
					if cookie == respCookie.Name {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Test %d sent invalid cookies: expected %s, got none", i+1, cookie)
				}
			}
		}
	}
}
