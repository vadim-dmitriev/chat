package http

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

// RegisterEndpoints регистрирует API Endtoint`ы, реализующие логику пакета auth
func RegisterEndpoints(auth auth.AuthServiceServer) {
	handler := handler{auth}

	http.HandleFunc("/api/v1/signup", handler.signUp)
	http.HandleFunc("/api/v1/signin", handler.signIn)
}
