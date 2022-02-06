package api

import (
	"encoding/json"
	"net/http"
)

var (
	DatabaseError = ApiError{Message: "internal database error", Code: http.StatusInternalServerError}
	InvalidJson   = ApiError{Message: "invalid json", Code: http.StatusUnprocessableEntity}
)

type ApiError struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func (ae ApiError) Send(w http.ResponseWriter) error {
	w.WriteHeader(ae.Code)
	return json.NewEncoder(w).Encode(ae)
}

var sessions = map[string]int{}

type article struct {
	ID       int    `json:"-"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Image    string `json:"image"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
	Link     string `json:"link"`
	Markdown string `json:"markdown"`
	Html     string `json:"html"`
}

func authorizedSession(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, false
	}
	for sessionId, authorId := range sessions {
		if sessionId == cookie.Value {
			return authorId, true
		}
	}
	return 0, false
}
