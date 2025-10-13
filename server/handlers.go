package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	CreateAt string `json:"created_at"`
}

// Response represents a standard API response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Sample users data
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", CreateAt: "2024-01-01"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", CreateAt: "2024-01-02"},
	{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", CreateAt: "2024-01-03"},
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go HTTP Server Demo</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 5px 0; border-radius: 5px; }
        .method { color: #007bff; font-weight: bold; }
    </style>
</head>
<body>
    <h1>🚀 Go HTTP Server Demo</h1>
    <p>Welcome to the Go HTTP server demonstration!</p>
    
    <h2>📋 Available Endpoints:</h2>
    <div class="endpoint">
        <span class="method">GET</span> / - Home page (this page)
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /health - Health check
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /time - Current time
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /users - List all users (HTML)
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /users/{id} - Get user by ID (HTML)
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /api/users - List all users (JSON)
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /api/users/{id} - Get user by ID (JSON)
    </div>
    
    <h2>🔗 Quick Links:</h2>
    <p><a href="/health">Health Check</a> | <a href="/time">Current Time</a> | <a href="/users">Users</a> | <a href="/api/users">API Users</a></p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: "Server is healthy",
		Data: map[string]any{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    "running",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TimeHandler handles time requests
func TimeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	response := Response{
		Success: true,
		Message: "Current time",
		Data: map[string]any{
			"time":      now.Format(time.RFC3339),
			"unix":      now.Unix(),
			"timezone":  now.Location().String(),
			"formatted": now.Format("2006-01-02 15:04:05"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UsersHandler handles HTML users list requests
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Users List</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>👥 Users List</h1>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Created At</th>
        </tr>`

	for _, user := range users {
		html += fmt.Sprintf(`
        <tr>
            <td>%d</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
        </tr>`, user.ID, user.Name, user.Email, user.CreateAt)
	}

	html += `
    </table>
    <p><a href="/">← Back to Home</a></p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// UserHandler handles individual user requests (HTML)
func UserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var foundUser *User
	for _, user := range users {
		if user.ID == id {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		http.NotFound(w, r)
		return
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>User %s</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .user-card { border: 1px solid #ddd; padding: 20px; border-radius: 5px; max-width: 400px; }
        .field { margin: 10px 0; }
        .label { font-weight: bold; }
    </style>
</head>
<body>
    <h1>👤 User Details</h1>
    <div class="user-card">
        <div class="field">
            <span class="label">ID:</span> %d
        </div>
        <div class="field">
            <span class="label">Name:</span> %s
        </div>
        <div class="field">
            <span class="label">Email:</span> %s
        </div>
        <div class="field">
            <span class="label">Created At:</span> %s
        </div>
    </div>
    <p><a href="/users">← Back to Users</a> | <a href="/">← Back to Home</a></p>
</body>
</html>`, foundUser.Name, foundUser.ID, foundUser.Name, foundUser.Email, foundUser.CreateAt)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// APIUsersHandler handles API users list requests (JSON)
func APIUsersHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// APIUserHandler handles individual user API requests (JSON)
func APIUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var foundUser *User
	for _, user := range users {
		if user.ID == id {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    foundUser,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
