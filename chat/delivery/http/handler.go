package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	authHTTP "github.com/vadim-dmitriev/chat/auth/delivery/http"

	"github.com/vadim-dmitriev/chat/chat"
)

type handler struct {
	chat chat.IChat
}

type chatResponseBody struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func (h handler) getConversations(w http.ResponseWriter, r *http.Request) {
	var encoder = json.NewEncoder(w)
	var allowedMethods = []string{
		http.MethodGet,
	}
	w.Header().Add("content-type", "application/json")

	if err := isMethodAllow(r.Method, allowedMethods); err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encodeResponse(encoder, nil, err)
		return
	}

	tokenCookie, err := r.Cookie(authHTTP.AuthHeaderName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		encodeResponse(encoder, nil, err)
		return
	}

	user, err := h.chat.GetUserByToken(tokenCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		encodeResponse(encoder, nil, err)
		return
	}

	conversations, err := h.chat.GetConversations(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encodeResponse(encoder, nil, err)
		return
	}

	encodeResponse(encoder, conversations, nil)
}

func (h handler) serachUser(w http.ResponseWriter, r *http.Request) {

}

func isMethodAllow(method string, allowedMethods []string) error {
	for _, allowedMethod := range allowedMethods {
		if method == allowedMethod {
			return nil
		}
	}

	return fmt.Errorf("method %s not allowed", method)
}

func encodeResponse(encoder *json.Encoder, data interface{}, err error) {
	var isSuccess = true
	var errorMessage = ""

	if err != nil {
		isSuccess = false
		errorMessage = err.Error()
	}

	encoder.Encode(
		chatResponseBody{
			Success: isSuccess,
			Data:    data,
			Error:   errorMessage,
		},
	)
}
