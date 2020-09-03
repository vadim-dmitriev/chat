package chat

import "github.com/vadim-dmitriev/chat/model"

type Chat struct {
	Repo IRepository
}

func (c Chat) SendMessage(message model.Message, from model.User, to model.Conversation) error {
	return nil
}

func (c Chat) GetUserByToken(token string) (model.User, error) {
	return c.Repo.GetUserByTokenChat(token)
}

func (c Chat) GetConversations(user model.User) ([]model.Conversation, error) {
	return c.Repo.GetConversations(user)
}
