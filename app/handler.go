package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (a App) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("session"); err == http.ErrNoCookie {
			http.Redirect(w, r, "/signin", http.StatusUnauthorized)
			return
		}

		newSessionCookie, _ := uuid.NewRandom()

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: newSessionCookie.String(),
			Path:  "/",
		})

		next(w, r)
	}
}

func (a App) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody = make(map[string]string)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		panic(err)
	}

	if !a.Storage.AuthUser(requestBody["username"], requestBody["password"]) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: requestBody["username"],
		Path:  "/",
	})

	newSessionCookie, _ := uuid.NewRandom()
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: newSessionCookie.String(),
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (a App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody = make(map[string]string)
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "ERROR while 'RegisterHandler': %s\n", err)
		return
	}
	if err := a.Storage.RegisterUser(requestBody["username"], requestBody["password"]); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR while 'RegisterHandler': %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a App) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrager := websocket.Upgrader{}
	conn, err := upgrager.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	username, _ := r.Cookie("username")

	a.Users[username.Value] = &User{username.Value, conn}

	go a.ServeUser(conn)
}
