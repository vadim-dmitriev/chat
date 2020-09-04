package chat

import "github.com/vadim-dmitriev/chat/model"

// IRepository интерфейс для хранения данных мессенджера
type IRepository interface {
	GetUserByTokenChat(token string) (model.User, error)

	SaveMessage(message model.Message, from model.User, to model.Conversation) error
	GetConversations(user model.User) ([]model.Conversation, error)
	GetMessages(conv model.Conversation, offset, limit int) ([]model.Message, error)
}
