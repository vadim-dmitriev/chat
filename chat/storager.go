package chat

import "github.com/vadim-dmitriev/chat/model"

// Storager интерфейс для хранения данных мессенджера
type Storager interface {
	SaveMessage(message model.Message, from model.User, to model.Conversation) error
}
