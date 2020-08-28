package auth

import (
	"fmt"

	"github.com/vadim-dmitriev/chat/model"
	"golang.org/x/crypto/bcrypt"
)

// JWT имплементирует интерфейс IAuth. Реализует аутентификацию
// пользователя при помощи JWT
type JWT struct {
	Repo IUserRepository
}

// SignIn аутентифицирует пользователя
func (a JWT) SignIn(username, password string) error {
	user, err := a.Repo.GetUser(username)
	if err != nil {
		return err
	}

	if !isPasswordMatch(password, user.Password) {
		return fmt.Errorf("password does not match")
	}

	return nil
}

// SignUp регистрирует нового пользователя
func (a JWT) SignUp(username, password string) error {
	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: encryptedPassword,
	}

	if err := a.Repo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func encryptPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), err
}

func isPasswordMatch(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	return true
}
