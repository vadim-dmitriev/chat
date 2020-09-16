package main

import (
	"net"

	"github.com/vadim-dmitriev/chat/auth"
	"github.com/vadim-dmitriev/chat/storage"
	"google.golang.org/grpc"
)

func main() {
	sqliteDB := storage.NewSqlite()

	a := auth.Session{
		Repo: sqliteDB,
	}

	srv := grpc.NewServer()

	auth.RegisterAuthServiceServer(srv, a)

	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	if err := srv.Serve(list); err != nil {
		panic(err)
	}

	// chat := chat.Chat{
	// 	Repo: sqliteDB,
	// }

	// authDeliveryHTTP.RegisterEndpoints(auth)
	// authDeliveryGRPC.RegisterGRPC(auth)
	// chatDeliveryWebsocket.RegisterUpgradeToWSEndpoint(chat)

	// server.RegisterHTTPStaticEndpoints(auth)

	// log.Println("Server started on 0.0.0.0:8080...")
	// defer log.Println("Server stopped")

	// log.Println(http.ListenAndServe("0.0.0.0:8080", nil))

}
