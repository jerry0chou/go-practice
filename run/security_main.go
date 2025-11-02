package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jerrychou/go-practice/security"
)

func main() {
	fmt.Println("=== Go Security Package Demo ===")

	// JWT Authentication Demo
	fmt.Println("\n1. JWT Authentication Demo")
	demoJWT()

	// OAuth Authentication Demo
	fmt.Println("\n2. OAuth Authentication Demo")
	demoOAuth()

	// RBAC Authorization Demo
	fmt.Println("\n3. RBAC Authorization Demo")
	demoRBAC()

	// Password Hashing Demo
	fmt.Println("\n4. Password Hashing Demo")
	demoPasswordHashing()

	// HTTPS/TLS Demo
	fmt.Println("\n5. HTTPS/TLS Demo")
	demoHTTPS()

	// Input Validation Demo
	fmt.Println("\n6. Input Validation Demo")
	demoInputValidation()
}

func demoJWT() {
	// Create JWT auth instance
	jwtAuth := security.NewJWTAuth("your-secret-key")

	// Generate token
	token, err := jwtAuth.GenerateToken("user123", "john_doe", []string{"user", "admin"}, 24)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return
	}

	fmt.Printf("Generated JWT Token: %s\n", token[:50]+"...")

	// Validate token
	claims, err := jwtAuth.ValidateToken(token)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		return
	}

	fmt.Printf("Token is valid for user: %s with roles: %v\n", claims.Username, claims.Roles)

	// Extract user info
	userID, username, roles, err := jwtAuth.ExtractUserInfo(token)
	if err != nil {
		log.Printf("Error extracting user info: %v", err)
		return
	}

	fmt.Printf("User Info - ID: %s, Username: %s, Roles: %v\n", userID, username, roles)
}

func demoOAuth() {
	// Create OAuth auth instance
	oauthAuth := security.NewOAuthAuth()

	// Add Google provider
	oauthAuth.AddProvider(security.GoogleProvider, security.OAuthConfig{
		ClientID:     "your-google-client-id",
		ClientSecret: "your-google-client-secret",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"openid", "profile", "email"},
	})

	// Generate auth URL
	authURL, err := oauthAuth.GetAuthURL(security.GoogleProvider, "random-state")
	if err != nil {
		log.Printf("Error generating auth URL: %v", err)
		return
	}

	fmt.Printf("OAuth Auth URL: %s\n", authURL)

	// Validate state parameter
	err = oauthAuth.ValidateState("random-state", "random-state")
	if err != nil {
		log.Printf("State validation error: %v", err)
		return
	}

	fmt.Println("State validation successful")
}

func demoRBAC() {
	// Create RBAC manager
	rbac := security.NewRBACManager()

	// Add permissions
	rbac.AddPermission(&security.Permission{
		Name:        "users:read",
		Resource:    "users",
		Action:      "read",
		Description: "Read user information",
	})

	rbac.AddPermission(&security.Permission{
		Name:        "users:write",
		Resource:    "users",
		Action:      "write",
		Description: "Write user information",
	})

	rbac.AddPermission(&security.Permission{
		Name:        "admin:all",
		Resource:    "*",
		Action:      "*",
		Description: "Full admin access",
	})

	// Add roles
	userRole := &security.Role{
		Name:        "user",
		Permissions: []string{"users:read"},
	}
	adminRole := &security.Role{
		Name:        "admin",
		Permissions: []string{"users:read", "users:write", "admin:all"},
	}

	rbac.AddRole(userRole)
	rbac.AddRole(adminRole)

	// Add users
	user := &security.User{
		ID:       "user1",
		Username: "john_doe",
		Email:    "john@example.com",
		Roles:    []string{"user"},
	}
	admin := &security.User{
		ID:       "admin1",
		Username: "admin_user",
		Email:    "admin@example.com",
		Roles:    []string{"admin"},
	}

	rbac.AddUser(user)
	rbac.AddUser(admin)

	// Test permissions
	fmt.Printf("User can read users: %v\n", rbac.HasPermission("user1", "users:read"))
	fmt.Printf("User can write users: %v\n", rbac.HasPermission("user1", "users:write"))
	fmt.Printf("Admin can write users: %v\n", rbac.HasPermission("admin1", "users:write"))
	fmt.Printf("Admin has admin role: %v\n", rbac.HasRole("admin1", "admin"))

	// Test resource access
	fmt.Printf("User can access users resource: %v\n", rbac.CheckResourceAccess("user1", "users", "read"))
	fmt.Printf("User can access admin resource: %v\n", rbac.CheckResourceAccess("user1", "admin", "write"))
}

func demoPasswordHashing() {
	// Create bcrypt hasher
	bcryptHasher := security.NewBcryptHasher(12)
	passwordManager := security.NewPasswordManager(bcryptHasher)

	// Hash password
	password := "SecurePassword123!"
	hash, err := passwordManager.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	fmt.Printf("Password hash: %s\n", hash[:50]+"...")

	// Verify password
	isValid := passwordManager.VerifyPassword(password, hash)
	fmt.Printf("Password verification: %v\n", isValid)

	// Test with wrong password
	isValid = passwordManager.VerifyPassword("wrongpassword", hash)
	fmt.Printf("Wrong password verification: %v\n", isValid)

	// Validate password strength
	result := passwordManager.ValidatePasswordStrength(password)
	if result != nil {
		fmt.Printf("Password strength validation error: %v\n", result)
	} else {
		fmt.Println("Password strength validation passed")
	}

	// Generate secure password
	securePassword, err := passwordManager.GenerateSecurePassword(16)
	if err != nil {
		log.Printf("Error generating secure password: %v", err)
		return
	}

	fmt.Printf("Generated secure password: %s\n", securePassword)
}

func demoHTTPS() {
	// Create TLS security instance
	tlsSecurity := security.NewTLSSecurity()

	// Set server config
	serverConfig := &security.TLSServerConfig{
		CertFile: "cert.pem",
		KeyFile:  "key.pem",
		MinTLS:   0x0303, // TLS 1.2
		MaxTLS:   0x0304, // TLS 1.3
	}
	tlsSecurity.SetServerConfig(serverConfig)

	// Create server TLS config
	serverTLSConfig, err := tlsSecurity.CreateServerTLSConfig()
	if err != nil {
		log.Printf("Error creating server TLS config: %v", err)
		return
	}

	fmt.Printf("Server TLS config created successfully\n")
	fmt.Printf("Min TLS version: %s\n", tlsSecurity.GetTLSVersionString(serverTLSConfig.MinVersion))
	fmt.Printf("Max TLS version: %s\n", tlsSecurity.GetTLSVersionString(serverTLSConfig.MaxVersion))

	// Create client TLS config
	_ = tlsSecurity.CreateClientTLSConfig()
	fmt.Printf("Client TLS config created successfully\n")

	// Create HTTPS client
	httpsClient := tlsSecurity.CreateHTTPSClient()
	fmt.Printf("HTTPS client created with timeout: %v\n", httpsClient.Timeout)

	// Generate self-signed certificate for development
	certPEM, keyPEM, err := tlsSecurity.GenerateSelfSignedCert("localhost")
	if err != nil {
		log.Printf("Error generating self-signed cert: %v", err)
		return
	}

	fmt.Printf("Self-signed certificate generated (length: %d bytes)\n", len(certPEM))
	fmt.Printf("Private key generated (length: %d bytes)\n", len(keyPEM))

	// Validate certificate
	err = tlsSecurity.ValidateCertificate(certPEM)
	if err != nil {
		log.Printf("Certificate validation error: %v", err)
		return
	}

	fmt.Println("Certificate validation passed")
}

func demoInputValidation() {
	// Create input validator
	validator := security.NewInputValidator()
	validator.CreateCommonRules()

	// Test data
	testInputs := map[string]string{
		"username": "john_doe123",
		"email":    "john@example.com",
		"password": "SecurePass123!",
		"name":     "John Doe",
		"url":      "https://example.com",
		"html":     "<script>alert('xss')</script><p>Safe content</p>",
		"sql":      "'; DROP TABLE users; --",
	}

	// Validate and sanitize inputs
	sanitized, errors := validator.ValidateAndSanitizeMap(testInputs)

	fmt.Println("Validation Results:")
	for field, value := range sanitized {
		fmt.Printf("%s: %s\n", field, value)
	}

	if len(errors) > 0 {
		fmt.Println("\nValidation Errors:")
		for _, err := range errors {
			fmt.Printf("- %s\n", err)
		}
	}

	// Test specific validations
	fmt.Println("\nSpecific Validation Tests:")

	// Email validation
	emailResult := validator.ValidateString("email", "invalid-email")
	fmt.Printf("Invalid email validation: Valid=%v, Errors=%v\n", emailResult.Valid, emailResult.Errors)

	// Password strength
	passwordResult := validator.ValidatePasswordStrength("weak")
	fmt.Printf("Weak password validation: Valid=%v, Errors=%v\n", passwordResult.Valid, passwordResult.Errors)

	// HTML sanitization
	sanitizedHTML := validator.SanitizeHTML(testInputs["html"])
	fmt.Printf("HTML sanitization: %s\n", sanitizedHTML)

	// SQL injection prevention
	sanitizedSQL := validator.PreventSQLInjection(testInputs["sql"])
	fmt.Printf("SQL injection prevention: %s\n", sanitizedSQL)

	// JSON validation
	jsonResult := validator.ValidateJSON(`{"name": "John", "age": 30}`)
	fmt.Printf("JSON validation: Valid=%v\n", jsonResult.Valid)
}

// HTTP handler for security headers demo
func securityHeadersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Security headers applied!"))
}

func demoSecurityHeaders() {
	// Create TLS security instance
	tlsSecurity := security.NewTLSSecurity()

	// Create HTTP server with security headers
	mux := http.NewServeMux()
	mux.HandleFunc("/", securityHeadersHandler)

	// Add security headers middleware
	handler := tlsSecurity.AddSecurityHeaders(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	fmt.Println("Starting server with security headers on :8080")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(1 * time.Second)

	// Test the server
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Printf("Error testing server: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Security headers in response:")
	for name, values := range resp.Header {
		if strings.HasPrefix(strings.ToLower(name), "x-") ||
			strings.ToLower(name) == "strict-transport-security" ||
			strings.ToLower(name) == "content-security-policy" ||
			strings.ToLower(name) == "referrer-policy" {
			fmt.Printf("%s: %v\n", name, values)
		}
	}

	server.Close()
}
