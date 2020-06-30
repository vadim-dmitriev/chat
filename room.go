package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

type room struct {
	clients map[*client]struct{}

	forward chan []byte
	join    chan *client
	leave   chan *client
}

func newRoom() *room {
	return &room{
		clients: make(map[*client]struct{}),
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
	}
}

func (r *room) Run() {
	for {
		select {
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}

		case newClient := <-r.join:
			r.clients[newClient] = struct{}{}

		case client := <-r.leave:
			delete(r.clients, client)
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalf("ServeHTTP: %s", err)
	}

	usernameCookie, err := req.Cookie("username")
	if err != nil {
		panic(err)
	}

	client := newClient(r, socket, usernameCookie.Value)
	defer func() {
		r.leave <- client
	}()

	r.join <- client

	go client.read()
	client.write()

}
