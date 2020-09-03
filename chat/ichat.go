package chat

import "github.com/vadim-dmitriev/chat/model"

// IChat интерфейс логики работы мессенджера
type IChat interface {
	GetUserByToken(token string) (model.User, error)

	SendMessage(message model.Message, from model.User, to model.Conversation) error
	GetConversations(user model.User) ([]model.Conversation, error)
}
