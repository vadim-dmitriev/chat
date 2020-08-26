package auth

import "github.com/vadim-dmitriev/chat/model"

// UserRepository интерфейс хранилища пользователей
type UserRepository interface {
	CreateUser(user model.User) error
	GetUser(username, password string) (model.User, error)
}
