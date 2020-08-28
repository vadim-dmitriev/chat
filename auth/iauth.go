package auth

import "github.com/vadim-dmitriev/chat/model"

// IAuth интерфейс логики, связанной с механизмом авторизации
type IAuth interface {
	SignUp(username, password string) error
	SignIn(username, password string) (token string, err error)
	ParseToken(token string) (model.User, error)
}
