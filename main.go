package main

import (
	"log"
	"net/http"
	"time"

	"github.com/vadim-dmitriev/chat/auth"
	authDeliveryHTTP "github.com/vadim-dmitriev/chat/auth/delivery/http"

	"github.com/vadim-dmitriev/chat/storage"
)

func main() {
	sqliteDB := storage.NewSqlite()

	auth := auth.JWT{
		Repo:         sqliteDB,
		Secret:       []byte("secret"),
		Method:       auth.SigningMethodHS256,
		ExpiringTime: time.Duration(1 * time.Hour),
	}
	authDeliveryHTTP.RegisterEndpoints(auth)

	log.Println("Server started on 0.0.0.0:8080...")
	defer log.Println("Server stopped")

	log.Println(http.ListenAndServe("0.0.0.0:8080", nil))
}
