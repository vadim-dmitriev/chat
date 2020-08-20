package app

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/vadim-dmitriev/chat/storage"
)

var (
	rpcNotFound = map[string]interface{}{
		"action": "not found",
	}
)

// App - связующая структура приложения
type App struct {
	Storage storage.Storager
	Users   map[string]*User
}

// New - конструктор структуры App
func New(s storage.Storager) App {
	return App{
		Storage: s,
		Users:   make(map[string]*User),
	}
}

// User описывает пользователя
type User struct {
	username string
	conn     *websocket.Conn
}

// ServeUser обслуживает соединение с пользователем
func (a App) ServeUser(conn *websocket.Conn) {
	request := make(map[string]interface{})
	var thisUsername string
	for usn, user := range a.Users {
		if conn == user.conn {
			thisUsername = usn
			break
		}
	}

	for {
		if err := conn.ReadJSON(&request); err != nil {
			fmt.Println(err)
			break
		}
		request["messageFrom"] = thisUsername

		if request["action"].(string) == "sendMessage" {
			for username, user := range a.Users {
				if username == request["conversationName"] {
					user.conn.WriteJSON(map[string]interface{}{
						"action": "newMessage",
						"value":  request["value"],
					})
					fmt.Println("sended to", username)
					err := a.Storage.SetMessage(request["message"].(string), request["messageFrom"].(string), request["conversationName"].(string))
					if err != nil {
						panic(err)
					}
					break
				}
			}
			continue
		}
		rpc := request["action"].(string)
		handler, ok := wsHandlers[rpc]
		if !ok {
			log.Println("not found")
			conn.WriteJSON(rpcNotFound)
			continue
		}
		log.Println(request["action"])

		conn.WriteJSON(
			handler(request, a.Storage),
		)
	}
}
