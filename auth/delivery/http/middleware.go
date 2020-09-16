package http

import (
	"fmt"
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
	"github.com/vadim-dmitriev/chat/model"
)

type Middleware struct {
	Auth auth.AuthServiceServer
}

func (m Middleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(authHeaderName)
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		token := model.Token{
			Value: authCookie.Value,
		}
		newToken, err := m.Auth.ParseToken(nil, &token)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; HttpOnly; Path=/", authHeaderName, newToken))
		next.ServeHTTP(w, r)
	}

}
