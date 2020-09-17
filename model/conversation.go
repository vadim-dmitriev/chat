package model

// Conversation модель беседы
type Conversation struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	IsDialog bool      `json:"is_dialog"`
	Messages []Message `json:"messages"`
}
