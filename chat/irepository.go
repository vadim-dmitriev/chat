package chat

import "github.com/vadim-dmitriev/chat/model"

// IRepository интерфейс для хранения данных мессенджера
type IRepository interface {
	SaveMessage(message model.Message, from model.User, to model.Conversation) error
}
