package app

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// WSHandler тип websocket обработчика (хэндлера)
type WSHandler func(*websocket.Conn)

var (
	wsHandlers = map[string]WSHandler{
		"searchUser": searUser(),
	}
)

func (a App) searchUser(conn *websocket.Conn) {
	fmt.Println(request["action"])

	response["action"] = "newConversationWith"
	response["username"] = request["username"]

	if a.Storage.IsUserExists(request["username"].(string)) {
		response["isUserExists"] = true
	} else {
		response["isUserExists"] = false
	}
	conn.WriteJSON(response)
}
