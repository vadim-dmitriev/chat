package main

import (
	"log"
	"net/http"

	"github.com/vadim-dmitriev/chat/server"
	"github.com/vadim-dmitriev/chat/storage"

	chatDeliveryHTTP "github.com/vadim-dmitriev/chat/auth/delivery/http"
	"github.com/vadim-dmitriev/chat/chat"
	chatDeliveryWebsocket "github.com/vadim-dmitriev/chat/chat/delivery/websocket"

	"github.com/vadim-dmitriev/chat/auth"
	authDeliveryHTTP "github.com/vadim-dmitriev/chat/auth/delivery/http"
)

func main() {
	sqliteDB := storage.NewSqlite()

	auth := auth.Session{
		Repo: sqliteDB,
	}
	authDeliveryHTTP.RegisterEndpoints(auth)

	chat := chat.Chat{
		Repo: sqliteDB,
	}
	chatDeliveryWebsocket.RegisterUpgradeToWSEndpoint(chat)
	chatDeliveryHTTP.RegisterEndpoints(chat)

	server.RegisterHTTPStaticEndpoints(auth)

	log.Println("Server started on 0.0.0.0:8080...")
	defer log.Println("Server stopped")

	log.Println(http.ListenAndServe("0.0.0.0:8080", nil))
}
