package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/mzzz-zzm/galleryblue/gen/go/users/v1/usersv1connect"
	"github.com/mzzz-zzm/galleryblue/internal/db"
	"github.com/mzzz-zzm/galleryblue/internal/handlers"
)

func main() {
	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	// Register AuthService handler
	authPath, authHandler := usersv1connect.NewAuthServiceHandler(&handlers.AuthServer{})
	mux.Handle(authPath, authHandler)

	// Register UserService handler
	userPath, userHandler := usersv1connect.NewUserServiceHandler(&handlers.UserServer{})
	mux.Handle(userPath, userHandler)

	// Register ImageService handler
	imagePath, imageHandler := usersv1connect.NewImageServiceHandler(&handlers.ImageServer{})
	mux.Handle(imagePath, imageHandler)

	// Add CORS support
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Connect-Protocol-Version",
			"X-User-ID",
			"Authorization",
		},
	}).Handler(mux)

	fmt.Println("Server executing on 0.0.0.0:8080")
	// Use h2c so we can serve HTTP/2 without TLS.
	if err := http.ListenAndServe(
		"0.0.0.0:8080",
		h2c.NewHandler(corsHandler, &http2.Server{}),
	); err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}
