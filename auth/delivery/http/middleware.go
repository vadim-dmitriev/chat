package http

import (
	"fmt"
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

type middleware struct {
	auth auth.IAuth
}

func (m middleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(authHeaderName)

		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := m.auth.ParseToken(token)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println(user)

		next(w, r)
	}

}
