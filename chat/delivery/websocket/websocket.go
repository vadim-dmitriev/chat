package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vadim-dmitriev/chat/chat"
)

type upgradeHandler struct {
	chat chat.IChat
}

func (uh upgradeHandler) upgrade(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	go serveConnection(conn, uh.chat)
}

func serveConnection(conn *websocket.Conn, chat chat.IChat)
