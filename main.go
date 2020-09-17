package main

import (
	"log"
	"net/http"

	chatDeliveryWebsocket "github.com/vadim-dmitriev/chat/chat/delivery/websocket"

	"github.com/vadim-dmitriev/chat/auth"
	authDeliveryHTTP "github.com/vadim-dmitriev/chat/auth/delivery/http"
	"github.com/vadim-dmitriev/chat/chat"
	"github.com/vadim-dmitriev/chat/server"

	"github.com/vadim-dmitriev/chat/storage"
)

func main() {
	sqliteDB := storage.NewSqlite()

	auth := auth.Session{
		Repo: sqliteDB,
	}

	chat := chat.Chat{
		Repo: sqliteDB,
	}

	authDeliveryHTTP.RegisterEndpoints(auth)
	chatDeliveryWebsocket.RegisterUpgradeToWSEndpoint(chat)

	server.RegisterHTTPStaticEndpoints(auth)

	log.Println("Server started on 0.0.0.0:8080...")
	defer log.Println("Server stopped")

	log.Println(http.ListenAndServe("0.0.0.0:8080", nil))
}
