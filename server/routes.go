package server

import (
	"net/http"
)

// SetupRoutes configures all the routes for the server
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Home page
	mux.HandleFunc("/", HomeHandler)

	// Health and utility endpoints
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/time", TimeHandler)

	// User endpoints (HTML)
	mux.HandleFunc("/users", UsersHandler)
	mux.HandleFunc("/users/", UserHandler)

	// API endpoints (JSON)
	mux.HandleFunc("/api/users", APIUsersHandler)
	mux.HandleFunc("/api/users/", APIUserHandler)

	// Static file serving (if needed)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	return mux
}

// SetupRoutesWithMiddleware configures routes with middleware
func SetupRoutesWithMiddleware() http.Handler {
	// Get the base routes
	handler := SetupRoutes()

	// Apply middleware in order (last applied is outermost)
	handler = SecurityMiddleware(handler)
	handler = CORSMiddleware(handler)
	handler = RateLimitMiddleware(handler)
	handler = LoggingMiddleware(handler)

	return handler
}
