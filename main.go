package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func serveMainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/html/chat.html")
}

func main() {
	room := newRoom()

	static := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", static))
	http.Handle("/ws", room)
	http.HandleFunc("/", serveMainPage)

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, syscall.SIGTERM)

	go room.Run()

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
