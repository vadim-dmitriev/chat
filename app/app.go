package app

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/vadim-dmitriev/chat/storage"
)

type App struct {
	Storage storage.Storager
	Users   map[*websocket.Conn]User
}

type User struct {
	conn *websocket.Conn
}

func (u User) ListenAndServe() {
	go u.read()
	u.write()
}

func (u User) read() {
	message := make(map[string]interface{})

	for {
		if err := u.conn.ReadJSON(&message); err != nil {
			fmt.Println(err)
			break
		}
		switch message["action"] {
		case "searchUser":
			fmt.Println("New SEARCHUSER message")
			u.conn.WriteJSON(map[string]interface{}{
				"action":      "newConversation",
				"name":        message["username"],
				"lastMessage": "Нет сообщений...",
			})
		}
	}
}

func (u User) write() {

}

func New(s storage.Storager) App {
	return App{
		Storage: s,
		Users:   make(map[*websocket.Conn]User),
	}
}
