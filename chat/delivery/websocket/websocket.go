package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vadim-dmitriev/chat/chat"
	"github.com/vadim-dmitriev/chat/model"
)

const (
	reqGetConversations = "getConversations"
	reqGetUser          = "getUser"
	reqNewMessage       = "newMessage"
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
	res := response{}
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
		switch req.Action {
		case reqGetUser:
			res.Action = reqGetUser
			res.Success = true
			res.Data = map[string]interface{}{
				"userID": c.ID,
			}
			conn.WriteJSON(res)

		case reqGetConversations:
			res.Action = reqGetConversations
			convs, err := chat.GetConversations(c.User)
			if err != nil {
				fmt.Println(err)
				res.Success = false
			} else {
				res.Success = true
			}
			res.Data = convs
			if err := conn.WriteJSON(res); err != nil {
				log.Printf("error while writing message to %s: %s", c.Name, err)
				continue
			}
		case reqNewMessage:
			res.Action = reqNewMessage
			reqData := req.Data.(map[string]interface{})
			conversationID := reqData["conversationID"].(string)
			messageText := reqData["text"].(string)
			msg := model.Message{
				From: c.User,
				To: model.Conversation{
					ID: conversationID,
				},
				Text: messageText,
			}

			err := chat.SendMessage(msg)

			if err != nil {
				res.Success = false
				res.Data = err.Error()
			} else {
				res.Success = true
				res.Data = msg
			}

			if err := conn.WriteJSON(res); err != nil {
				panic(err)
			}
		}
	}
}
