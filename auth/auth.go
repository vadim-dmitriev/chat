package auth

import (
	"fmt"

	"github.com/vadim-dmitriev/chat/model"
	"golang.org/x/crypto/bcrypt"
)

// Auth имплементирует интерфейс IAuth
type Auth struct {
	repo UserRepository
}

// SignIn аутентифицирует пользователя
func (a Auth) SignIn(username, password string) error {
	user, err := a.repo.GetUser(username)
	if err != nil {
		return err
	}

	if !isPasswordMatch(password, user.Password) {
		return fmt.Errorf("password does not match")
	}

	return nil
}

// SignUp регистрирует нового пользователя
func (a Auth) SignUp(username, password string) error {
	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: encryptedPassword,
	}

	if err := a.repo.CreateUser(user); err != nil {
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
