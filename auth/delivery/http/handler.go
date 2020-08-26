package http

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

type handler struct {
	auth auth.IAuth
}

func newHandler(auth auth.IAuth) Handler {
	return Handler{
		auth,
	}
}

func (h Handler) signUp(w http.ResponseWriter, r *http.Request) {}

func (h Handler) signIn(w http.ResponseWriter, r *http.Request) {}
