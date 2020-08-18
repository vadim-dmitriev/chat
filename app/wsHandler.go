package app

import (
	"github.com/vadim-dmitriev/chat/storage"
)

// WSHandler тип websocket обработчика (хэндлера)
type WSHandler func(map[string]interface{}, storage.Storager) map[string]interface{}

var (
	wsHandlers = map[string]WSHandler{
		"searchUser": searchUserWSHandler,
	}
)

func searchUserWSHandler(request map[string]interface{}, s storage.Storager) map[string]interface{} {
	var response = make(map[string]interface{}, 3)
	response["action"] = "searchUser"
	response["isUserExists"] = false

	username := request["username"].(string)
	if s.IsUserExists(username) {
		response["isUserExists"] = true
		response["newConversationWith"] = username
	}

	return response
}
