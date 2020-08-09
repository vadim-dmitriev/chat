package app

import (
	"github.com/vadim-dmitriev/chat/storage"
)

type App struct {
	Storage storage.Storager
}

func New(s storage.Storager) App {
	return App{
		Storage: s,
	}
}
