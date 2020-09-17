package http

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/chat"
)

func RegisterEndpoints(chat chat.IChat) {
	handler := handler{chat}

	http.HandleFunc("/api/v1/conversations", handler.getConversations)
	http.HandleFunc("/api/v1/users/search", handler.serachUser)
}
