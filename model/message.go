package model

import "time"

// Message модель сообщения
type Message struct {
	Text     string
	Datetime time.Time
}
