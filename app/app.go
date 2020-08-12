package app

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/vadim-dmitriev/chat/storage"
)

type App struct {
	Storage storage.Storager
	Users   map[string]*User
}

type User struct {
	username string
	conn     *websocket.Conn
}

func (a App) ServeUser(conn *websocket.Conn) {
	// read function
	go func() {
		request := make(map[string]interface{})
		response := make(map[string]interface{})

		for {
			if err := conn.ReadJSON(&request); err != nil {
				fmt.Println(err)
				break
			}
			switch request["action"] {
			case "searchUser":
				fmt.Println(request["action"])

				response["action"] = "newConversationWith"
				response["username"] = request["username"]

				if a.Storage.IsUserExists(request["username"].(string)) {
					response["isUserExists"] = true
				} else {
					response["isUserExists"] = false
				}
				conn.WriteJSON(response)

			case "sendMessage":
				fmt.Println(request["action"], request)
			}
		}
	}()

	// write function
	// TODO
}

func read() {

}

func New(s storage.Storager) App {
	return App{
		Storage: s,
		Users:   make(map[string]*User),
	}
}
