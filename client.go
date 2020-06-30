package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	name []byte
	room *room
	conn *websocket.Conn
	send chan []byte
}

func newClient(r *room, conn *websocket.Conn, name string) *client {
	return &client{
		name: []byte(name),
		room: r,
		conn: conn,
		send: make(chan []byte),
	}
}

func (c *client) read() {
	defer c.conn.Close()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		c.room.forward <- append(append(c.name, ':', ' '), msg...)
	}

}

func (c *client) write() {
	defer c.conn.Close()

	for {
		for msg := range c.send {
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}

}
