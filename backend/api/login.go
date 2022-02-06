package api

import (
	"TheAdversary/database"
	"TheAdversary/schema"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
)

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		InvalidJson.Send(w)
		return
	}
	var author schema.Author
	database.GetDB().Table("author").Select("id", "password").Where("name = ?", payload.Username).Take(&author)
	if author.Id == 0 || bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(payload.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID := sessionId()

	sessions[sessionID] = author.Id

	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})
	w.WriteHeader(http.StatusOK)
}

func sessionId() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 32)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
