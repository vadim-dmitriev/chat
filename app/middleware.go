package app

import (
	"net/http"

	"github.com/google/uuid"
)

// AuthMiddleware проверяет авторизован ли пользователь
func (a App) AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Проверка на "Есть ли у клиента кука сессии?"
		currentSessionCookie, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/signin", http.StatusUnauthorized)
			return
		}
		username, err := a.Storage.GetUsernameByCookie(currentSessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusUnauthorized)
			return
		}

		newSessionCookie, _ := uuid.NewRandom()
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: newSessionCookie.String(),
			Path:  "/",
		})

		if err := a.Storage.UpdateUserSessionCookie(newSessionCookie.String(), username); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})

}
