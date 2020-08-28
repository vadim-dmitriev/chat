package auth

import "github.com/vadim-dmitriev/chat/model"

// IUserRepository интерфейс хранилища пользователей
type IUserRepository interface {
	CreateUser(user model.User) error
	GetUser(username string) (model.User, error)
}
