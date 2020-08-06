package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func serveMainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/html/login.html")
}

func auth(w http.ResponseWriter, r *http.Request) {
	var requestBody = make(map[string]string)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	if requestBody["login"] == "vadim" && requestBody["password"] == "1" {
		fmt.Fprintf(w, "%s", "Success")
	} else {
		fmt.Fprintf(w, "%s", "Login or password incorrect")
	}
}

func main() {
	room := newRoom()

	static := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", static))
	http.Handle("/ws", room)
	http.HandleFunc("/", serveMainPage)
	http.HandleFunc("/api/v1/auth", auth)

	// http.HandleFunc("/", serveMainPage)

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
