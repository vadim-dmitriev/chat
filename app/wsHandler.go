package app

import (
	"github.com/vadim-dmitriev/chat/storage"
)

// WSHandler тип websocket обработчика (хэндлера)
type WSHandler func(map[string]interface{}, storage.Storager) map[string]interface{}

var (
	wsHandlers = map[string]WSHandler{
		"searchUser":       searchUserWSHandler,
		"getConversations": getConversations,
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

func getConversations(request map[string]interface{}, s storage.Storager) map[string]interface{} {
	var response = make(map[string]interface{}, 2)
	response["action"] = "conversations"

	conversations, err := s.GetUserConversations(request["username"].(string))
	if err != nil {
		response["success"] = false
		response["error"] = err.Error()
		return response
	}

	response["conversations"] = conversations

	return response
}
