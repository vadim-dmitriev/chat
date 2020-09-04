package model

import "time"

// Message модель сообщения
type Message struct {
	ID       string       `json:"id"`
	Text     string       `json:"text"`
	From     User         `json:"from"`
	To       Conversation `json:"to"`
	Datetime time.Time    `json:"datetime"`
}
