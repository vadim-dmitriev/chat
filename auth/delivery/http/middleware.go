package http

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

type middleware struct {
	auth auth.IAuth
}

func (m middleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}

}
