package http

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

// RegisterEndpoints регистрирует API Endtoint`ы, реализующие логику пакета auth
func RegisterEndpoints(auth auth.IAuth) {
	handler := handler{auth}
	mw := middleware{auth}
	http.HandleFunc("/api/v1/signup", handler.signUp)
	// http.HandleFunc("/api/v1/signin", handler.signIn)
	http.HandleFunc("/api/v1/signin", mw.Handle(handler.signIn))
}
