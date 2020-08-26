package chat

import "github.com/vadim-dmitriev/chat/model"

// Repository интерфейс для хранения данных мессенджера
type Repository interface {
	SaveMessage(message model.Message, from model.User, to model.Conversation) error
}
