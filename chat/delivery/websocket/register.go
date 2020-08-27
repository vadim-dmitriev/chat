package websocket

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/chat"
)

// RegisterUpgradeToWSEndpoint регистрирует API Endtoint, на котором
// происходит upgrade http соединения до WebSocket соединения
func RegisterUpgradeToWSEndpoint(chat chat.IChat) {
	upgraderHandler := upgradeHandler{chat}

	http.HandleFunc("/api/v1/ws", upgraderHandler.upgrade)
}
