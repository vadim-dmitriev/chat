package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type authHandler struct {
}

func newAuthHandler() authHandler {
	return authHandler{}
}

func (ah authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestBody = make(map[string]string)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		panic(err)
	}

	fmt.Println(requestBody)

	if requestBody["login"] == "vadim" && requestBody["password"] == "1" {
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: requestBody["login"],
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		// http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
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

type registerHandler struct{}

func newRegisterHandler() registerHandler {
	return registerHandler{}
}

func (rh registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type chatHandler struct{}

func newChatHandler() chatHandler {
	return chatHandler{}
}

func (ch chatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("username")
	if err == http.ErrNoCookie {
		fmt.Println(err)
		http.Redirect(w, r, "/singin", http.StatusMovedPermanently)
		return
	}
	http.ServeFile(w, r, "static/html/chat.html")

}
