package main

import (
	"net/http"
)

type soundNotification struct {
}

func newSound() *soundNotification {
	return &soundNotification{}
}

func (sn *soundNotification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "swiftly.ogg")
}
