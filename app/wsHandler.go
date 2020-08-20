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
		"sendMessage":      sendMessage,
	}
)

func searchUserWSHandler(request map[string]interface{}, s storage.Storager) map[string]interface{} {
	var response = make(map[string]interface{}, 3)
	response["action"] = "searchUser"
	response["isUserExists"] = false

	username := request["username"].(string)
	if username == request["messageFrom"].(string) {
		return response
	}

	if s.IsUserExists(username) {
		response["isUserExists"] = true
		response["newConversationWith"] = username
	}

	return response
}

func getConversations(request map[string]interface{}, s storage.Storager) map[string]interface{} {
	var response = make(map[string]interface{})
	response["action"] = "conversations"

	conversations, err := s.GetUserConversations(request["messageFrom"].(string))
	if err != nil {
		response["success"] = false
		response["error"] = err.Error()
		return response
	}
	response["conversations"] = conversations

	return response
}

func sendMessage(request map[string]interface{}, s storage.Storager) map[string]interface{} {
	var response = make(map[string]interface{})
	response["success"] = true
	// map[action:sendMessage conversationName:1 message:1 messageFrom:vadim]
	err := s.SetMessage(request["message"].(string), request["messageFrom"].(string), request["conversationName"].(string))
	if err != nil {
		response["success"] = false
	}

	return response
}
