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

func New(s storage.Storager) App {
	return App{
		AuthHandler:      newAuthHandler(s),
		ChatHandler:      newChatHandler(),
		RegisterHandler:  newRegisterHandler(s),
		WebSocketHandler: newWebSocketHandler(),
		Storage:          storage.NewSqlite(),
	}
}
