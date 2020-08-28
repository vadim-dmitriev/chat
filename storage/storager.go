package storage

import (
	"github.com/vadim-dmitriev/chat/auth"
	"github.com/vadim-dmitriev/chat/chat"
)

type Storager interface {
	auth.UserRepository
	chat.Repository
}
