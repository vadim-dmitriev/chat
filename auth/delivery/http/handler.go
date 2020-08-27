package http

import (
	"encoding/json"
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

type handler struct {
	auth auth.IAuth
}

type authBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := authBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.auth.SignUp(body.Username, body.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
}

func (h handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := authBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.auth.SignIn(body.Username, body.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
