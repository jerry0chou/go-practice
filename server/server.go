package server

import (
	"fmt"
	"net/http"
	"time"
)

// Server represents the HTTP server configuration
type Server struct {
	Port    string
	Handler http.Handler
}

// New creates a new server instance
func New(port string) *Server {
	return &Server{
		Port: port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	server := &http.Server{
		Addr:         ":" + s.Port,
		Handler:      s.Handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("🚀 HTTP Server starting on port %s\n", s.Port)
	fmt.Printf("📋 Available endpoints:\n")
	fmt.Printf("   GET  /           - Home page\n")
	fmt.Printf("   GET  /health     - Health check\n")
	fmt.Printf("   GET  /time       - Current time\n")
	fmt.Printf("   GET  /users      - List all users\n")
	fmt.Printf("   GET  /users/{id} - Get user by ID\n")
	fmt.Printf("   GET  /api/users  - API: List all users (JSON)\n")
	fmt.Printf("   GET  /api/users/{id} - API: Get user by ID (JSON)\n")

	return server.ListenAndServe()
}

// SetHandler sets the HTTP handler for the server
func (s *Server) SetHandler(handler http.Handler) {
	s.Handler = handler
}
