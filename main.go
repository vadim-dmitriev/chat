package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

type templateHandler struct {
	sync.Once

	filename string
	template *template.Template
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.Do(func() {
		th.template = template.Must(template.ParseFiles(path.Join("templates", th.filename)))
	})

	if err := th.template.Execute(w, r); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}

func main() {
	room := newRoom()
	sound := newSound()

	http.Handle("/notify", sound)
	http.Handle("/room", room)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/", &templateHandler{filename: "mainPage.html"})

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, syscall.SIGTERM)

	go room.Run()

	server := http.Server{
		Addr: ":8081",
	}
	fmt.Println("Server started on :8081...")
	go server.ListenAndServe()

	fmt.Printf("\rServer shuting down: %s\n", <-exitChan)
	if err := server.Shutdown(nil); err != nil {
		panic(err)
	}
}
