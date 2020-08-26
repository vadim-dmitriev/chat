package auth

// Auth имплементирует интерфейс IAuth
type Auth struct{}

// SignIn аутентифицирует пользователя
func (a Auth) SignIn(username, password string) error {
	return nil
}

// SignUp регистрирует нового пользователя
func (a Auth) SignUp(username, password string) error {
	return nil
}
