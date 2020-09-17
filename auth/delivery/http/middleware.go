package http

import (
	"fmt"
	"net/http"

	"github.com/vadim-dmitriev/chat/auth"
)

type Middleware struct {
	Auth auth.IAuth
}

func (m Middleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(AuthHeaderName)
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		_, newToken, err := m.Auth.ParseToken(authCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; HttpOnly; Path=/", AuthHeaderName, newToken))
		next.ServeHTTP(w, r)
	}

}
