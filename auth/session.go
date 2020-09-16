package auth

import (
	context "context"

	model "github.com/vadim-dmitriev/chat/model"
)

type Session struct {
	Repo IUserRepository
}

func (s Session) SignIn(context.Context, *model.User) (*model.Token, error) {
	return nil, nil
}

func (s Session) SignUp(context.Context, *model.User) (*Empty, error) {
	return nil, nil
}

func (s Session) ParseToken(context.Context, *model.Token) (*model.Token, error) {
	return nil, nil
}
