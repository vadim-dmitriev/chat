package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

type handler struct {
	auth auth.IAuth
}

type authRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponseBody struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (h handler) signUp(w http.ResponseWriter, r *http.Request) {
	var encoder = json.NewEncoder(w)

	if r.Method != http.MethodPost {
		err := fmt.Errorf("method %s not allowed", r.Method)
		encodeResponse(encoder, err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := authRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		encodeResponse(encoder, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.auth.SignUp(body.Username, body.Password); err != nil {
		encodeResponse(encoder, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
}

func (h handler) signIn(w http.ResponseWriter, r *http.Request) {
	var encoder = json.NewEncoder(w)

	if r.Method != http.MethodPost {
		err := fmt.Errorf("method %s not allowed", r.Method)
		encodeResponse(encoder, err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := authRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		encodeResponse(encoder, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.auth.SignIn(body.Username, body.Password); err != nil {
		encodeResponse(encoder, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func encodeResponse(encoder *json.Encoder, err error) {
	var isSuccess = true
	var errorMessage = ""

	if err != nil {
		isSuccess = false
		errorMessage = err.Error()
	}

	encoder.Encode(
		authResponseBody{
			Success: isSuccess,
			Error:   errorMessage,
		},
	)
}
