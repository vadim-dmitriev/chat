package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vadim-dmitriev/chat/storage"
)

type authHandler struct {
	s storage.Storager
}

func newAuthHandler(s storage.Storager) authHandler {
	return authHandler{
		s: s,
	}
}

func (ah authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestBody = make(map[string]string)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		panic(err)
	}

	if ah.s.AuthUser(requestBody["login"], requestBody["password"]) {
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: requestBody["login"],
			Path:  "/",
		})
		w.Write([]byte("ok"))
	}
}

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func newWebSocketHandler() webSocketHandler {
	upgrader := websocket.Upgrader{}

	return webSocketHandler{
		upgrader: upgrader,
	}
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("ServeHTTP: %s", err)
	}
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

}

type registerHandler struct {
	s storage.Storager
}

func newRegisterHandler(s storage.Storager) registerHandler {
	return registerHandler{
		s,
	}
}

func (rh registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestBody = make(map[string]string)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		panic(err)
	}

	fmt.Println(requestBody)

	if err := rh.s.RegisterUser(requestBody["login"], requestBody["password"]); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "ok")

}

type chatHandler struct{}

func newChatHandler() chatHandler {
	return chatHandler{}
}

func (ch chatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("username")
	if err == http.ErrNoCookie {
		fmt.Println("NO COOKIES")
		http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
		return
	}

	http.ServeFile(w, r, "static/html/chat.html")
}
