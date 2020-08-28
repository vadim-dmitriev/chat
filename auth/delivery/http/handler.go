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
	var allowedMethods = []string{
		http.MethodPost,
	}
	w.Header().Add("content-type", "application/json")

	if err := isMethodAllow(r.Method, allowedMethods); err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encodeResponse(encoder, err)
		return
	}

	body := authRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encodeResponse(encoder, err)
		return
	}

	if err := h.auth.SignUp(body.Username, body.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encodeResponse(encoder, err)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
}

func (h handler) signIn(w http.ResponseWriter, r *http.Request) {
	var encoder = json.NewEncoder(w)
	var allowedMethods = []string{
		http.MethodPost,
	}
	w.Header().Add("content-type", "application/json")

	if err := isMethodAllow(r.Method, allowedMethods); err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encodeResponse(encoder, err)
		return
	}

	body := authRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encodeResponse(encoder, err)
		return
	}

	if err := h.auth.SignIn(body.Username, body.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		encodeResponse(encoder, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func isMethodAllow(method string, allowedMethods []string) error {
	for _, allowedMethod := range allowedMethods {
		if method == allowedMethod {
			return nil
		}
	}

	return fmt.Errorf("method %s not allowed", method)
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
