package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	usersv1 "github.com/mzzz-zzm/galleryblue/gen/go/users/v1"
	"github.com/mzzz-zzm/galleryblue/gen/go/users/v1/usersv1connect"
)

type UserServer struct{}

func (s *UserServer) GetUser(
	ctx context.Context,
	req *connect.Request[usersv1.GetUserRequest],
) (*connect.Response[usersv1.GetUserResponse], error) {
	log.Printf("Request headers: %v", req.Header())
	res := connect.NewResponse(&usersv1.GetUserResponse{
		Id:    "123",
		Name:  "Jane Doe",
		Email: "jane@example.com",
	})
	return res, nil
}

func main() {
	mux := http.NewServeMux()
	path, handler := usersv1connect.NewUserServiceHandler(&UserServer{})
	mux.Handle(path, handler)

	// Add CORS support
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Vite default port
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Connect-Protocol-Version"},
	}).Handler(mux)

	fmt.Println("Server executing on localhost:8080")
	// Use h2c so we can serve HTTP/2 without TLS.
	if err := http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(corsHandler, &http2.Server{}),
	); err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}
