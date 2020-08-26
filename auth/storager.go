package auth

import "github.com/vadim-dmitriev/chat/model"

// Storager интерфейс хранилища пользователей
type Storager interface {
	CreateUser(user model.User) error
	GetUser(username, password string) (model.User, error)
}
