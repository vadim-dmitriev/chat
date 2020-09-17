package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vadim-dmitriev/chat/chat"
	"github.com/vadim-dmitriev/chat/model"
)

type upgradeHandler struct {
	chat chat.IChat
}

type client struct {
	model.User
}

type request struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type response struct {
	Action  string      `json:"action"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func (uh upgradeHandler) upgrade(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("Authorization")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	user, err := uh.chat.GetUserByToken(tokenCookie.Value)

	go client{user}.Serve(conn, uh.chat)
	log.Printf("client %s connected", user.Name)

}

func (c client) Serve(conn *websocket.Conn, chat chat.IChat) {
	defer func() {
		conn.WriteMessage(websocket.CloseMessage, nil)
		conn.Close()
		log.Printf("client %s disconnected", c.Name)
	}()

	req := request{}

	for {
		err := conn.ReadJSON(&req)
		switch err.(type) {
		case nil:
			break
		case *websocket.CloseError:
			return
		default:
			log.Printf("error while reading %s message: %s", c.Name, err)
			continue
		}
		log.Printf("message from client %s: %s", c.Name, req.Action)
	}
}
