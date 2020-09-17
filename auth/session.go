package auth

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vadim-dmitriev/chat/model"
)

type Session struct {
	Repo IUserRepository
}

func (s Session) SignIn(username, password string) (token string, err error) {
	user, err := s.Repo.GetUser(username)
	if err != nil {
		return "", err
	}

	if !isPasswordMatch(password, user.Password) {
		return "", fmt.Errorf("password does not match")
	}

	t, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	if err := s.Repo.SetUserToken(user, t.String()); err != nil {
		return "", fmt.Errorf("internal error")
	}

	return t.String(), nil
}

func (s Session) SignUp(username, password string) error {
	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     username,
		Password: encryptedPassword,
	}

	return s.Repo.CreateUser(user)
}

func (s Session) ParseToken(token string) (model.User, string, error) {
	user, err := s.Repo.GetUserByToken(token)
	if err != nil {
		return model.User{}, "", err
	}

	newToken, err := uuid.NewRandom()
	if err != nil {
		return model.User{}, "", err
	}

	if err := s.Repo.SetUserToken(user, newToken.String()); err != nil {
		return model.User{}, "", err
	}

	return user, newToken.String(), nil
}
