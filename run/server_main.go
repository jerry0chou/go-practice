package main

import (
	"log"
	"os"

	"github.com/jerrychou/go-practice/server"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new server instance
	srv := server.New(port)

	// Setup routes with middleware
	handler := server.SetupRoutesWithMiddleware()
	srv.SetHandler(handler)

	// Start the server
	log.Printf("Starting server on port %s", port)
	if err := srv.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
