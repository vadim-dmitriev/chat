package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	authDeliveryHTTP "github.com/vadim-dmitriev/chat/auth/delivery/http"

	"github.com/vadim-dmitriev/chat/auth"
	"github.com/vadim-dmitriev/chat/storage"
)

func main() {
	sqliteDB := storage.NewSqlite()

	auth := auth.Auth{
		Repo: sqliteDB,
	}

	authDeliveryHTTP.RegisterEndpoints(auth)

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, syscall.SIGTERM)

	server := http.Server{
		Addr: "0.0.0.0:8080",
	}
	fmt.Println("Server started on 0.0.0.0:8080...")
	go server.ListenAndServe()

	fmt.Printf("\rServer shuting down: %s\n", <-exitChan)
	if err := server.Shutdown(nil); err != nil {
		panic(err)
	}
}
