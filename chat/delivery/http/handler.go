package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vadim-dmitriev/chat/chat"
)

type handler struct {
	chat chat.IChat
}

type chatResponseBody struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (h handler) getConversations(w http.ResponseWriter, r *http.Request) {

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

func encodeResponse(encoder *json.Encoder, err error) {
	var isSuccess = true
	var errorMessage = ""

	if err != nil {
		isSuccess = false
		errorMessage = err.Error()
	}

	encoder.Encode(
		chatResponseBody{
			Success: isSuccess,
			Error:   errorMessage,
		},
	)
}
