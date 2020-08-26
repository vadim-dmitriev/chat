package auth

// IAuth интерфейс логики, связанной с механизмом авторизации
type IAuth interface {
	SignUp(username, password string) error
	SignIn(username, password string) error
}
