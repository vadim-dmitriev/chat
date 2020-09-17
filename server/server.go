package server

import (
	"net/http"

	authDeliveryHTTP "github.com/vadim-dmitriev/chat/auth/delivery/http"

	"github.com/vadim-dmitriev/chat/auth"
)

func RegisterHTTPStaticEndpoints(auth auth.IAuth) {
	signInTemplate, err := newTemplate("./server/static/html/signin.html")
	signUpTemplate, err := newTemplate("./server/static/html/signup.html")
	chatTemplate, err := newTemplate("./server/static/html/chat.html")
	if err != nil {
		panic(err)
	}

	authMiddleware := authDeliveryHTTP.Middleware{auth}

	http.Handle("/signin", signInTemplate)
	http.Handle("/signup", signUpTemplate)
	http.Handle("/", authMiddleware.Handle(chatTemplate.ServeHTTP))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))
}
