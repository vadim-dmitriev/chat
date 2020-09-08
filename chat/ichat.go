package chat

import "github.com/vadim-dmitriev/chat/model"

// IChat интерфейс логики работы мессенджера
type IChat interface {
	GetUserByToken(token string) (model.User, error)
	NewUser(user model.User, conn interface{})

	SendMessage(message model.Message) error
	GetConversations(user model.User) ([]model.Conversation, error)
}
