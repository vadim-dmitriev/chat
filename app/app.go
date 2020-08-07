package app

import (
	"net/http"

	"github.com/vadim-dmitriev/chat/storage"
)

type App struct {
	AuthHandler      http.Handler
	ChatHandler      http.Handler
	RegisterHandler  http.Handler
	WebSocketHandler http.Handler

	Storage storage.Storager
}

func New() App {
	return App{
		AuthHandler:      newAuthHandler(),
		ChatHandler:      newChatHandler(),
		RegisterHandler:  newRegisterHandler(),
		WebSocketHandler: newWebSocketHandler(),
		Storage:          storage.NewSqlite(),
	}
}
