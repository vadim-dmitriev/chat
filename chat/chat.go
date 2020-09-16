package chat

import (
	"github.com/vadim-dmitriev/chat/model"
)

type Chat struct {
	Repo IRepository
	// Users map[model.User]*websocket.Conn
}

func (c Chat) NewUser(user model.User, conn interface{}) {
	// c.Users[user] = conn.(*websocket.Conn)
}

func (c Chat) SendMessage(message model.Message) error {
	// TODO: Send message to user(s)
	return c.Repo.SaveMessage(message)
}

func (c Chat) GetUserByToken(token string) (model.User, error) {
	return c.Repo.GetUserByTokenChat(token)
}

func (c Chat) GetConversations(user model.User) ([]model.Conversation, error) {
	return c.Repo.GetConversations(user)
}
