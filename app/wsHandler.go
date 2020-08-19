package app

import (
	"fmt"

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

	// convResponse := make([]map[string]interface{}, len(conversations), len(conversations))
	// for i, convName := range conversations {
	// 	conversation := make(map[string]interface{})
	// 	conversation["name"] = convName
	// 	convResponse[i] = conversation
	// }
	fmt.Println(conversations)
	response["conversations"] = conversations

	return response
}

func sendMessage(request map[string]interface{}, s storage.Storager) map[string]interface{} {
	var response = make(map[string]interface{})
	fmt.Println(request)
	return response
}
