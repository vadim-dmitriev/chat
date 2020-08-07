package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vadim-dmitriev/chat/app"
)

func main() {
	app := app.New()

	http.Handle("/api/v1/auth", app.AuthHandler)
	http.HandleFunc("/singin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/login.html")
	})

	http.Handle("/singup", app.RegisterHandler)
	http.Handle("/ws", app.WebSocketHandler)
	http.Handle("/", app.ChatHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, syscall.SIGTERM)

	server := http.Server{
		Addr: "0.0.0.0:8081",
	}
	fmt.Println("Server started on 0.0.0.0:8081...")
	go server.ListenAndServe()

	fmt.Printf("\rServer shuting down: %s\n", <-exitChan)
	if err := server.Shutdown(nil); err != nil {
		panic(err)
	}
}
