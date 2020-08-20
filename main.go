package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vadim-dmitriev/chat/app"
	"github.com/vadim-dmitriev/chat/storage"
)

func main() {
	s := storage.NewSqlite()
	app := app.New(s)

	// API Paths
	http.HandleFunc("/api/v1/auth", app.AuthHandler)
	http.HandleFunc("/api/v1/register", app.RegisterHandler)
	http.HandleFunc("/api/v1/ws", app.WebSocketHandler)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	// Pages
	http.Handle("/signin", page("static/html/signin.html"))
	http.Handle("/signup", page("static/html/signup.html"))
	http.Handle("/", app.AuthMiddleware(
		page("static/html/chat.html"),
	))

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
