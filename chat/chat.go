package chat

import "github.com/vadim-dmitriev/chat/model"

type chat struct {
}

func (c chat) SendMessage(message model.Message, from model.User, to model.Conversation) error {
	return nil
}
