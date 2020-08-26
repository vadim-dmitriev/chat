package chat

import "github.com/vadim-dmitriev/chat/model"

// IChat интерфейс логики работы мессенджера
type IChat interface {
	SendMessage(message model.Message, from model.User, to model.Conversation) error
}
