package http

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

// RegisterHTTPEndpoints регистрирует API Endtoint`ы, реализующие логику пакета auth
func RegisterHTTPEndpoints(auth auth.IAuth) {
	handler := handler{auth}

	http.HandleFunc("/api/v1/signup", handler.signUp)
	http.HandleFunc("/api/v1/signin", handler.signIn)
}
