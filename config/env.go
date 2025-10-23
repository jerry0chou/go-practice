package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvConfig handles environment variable configuration following 12-factor app principles
type EnvConfig struct {
	// Application settings
	AppName        string
	AppVersion     string
	AppEnvironment string // development, staging, production

	// Server settings
	ServerHost string
	ServerPort int

	// Database settings
	DatabaseURL      string
	DatabaseMaxConns int
	DatabaseTimeout  time.Duration

	// Logging settings
	LogLevel  string
	LogFormat string

	// Feature flags
	EnableMetrics bool
	EnableDebug   bool
	EnableCORS    bool

	// External services
	RedisURL     string
	CacheTimeout time.Duration

	// Security
	JWTSecret     string
	SessionSecret string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*EnvConfig, error) {
	config := &EnvConfig{}

	// Application settings
	config.AppName = getEnv("APP_NAME", "go-practice")
	config.AppVersion = getEnv("APP_VERSION", "1.0.0")
	config.AppEnvironment = getEnv("APP_ENV", "development")

	// Server settings
	config.ServerHost = getEnv("SERVER_HOST", "localhost")
	config.ServerPort = getEnvAsInt("SERVER_PORT", 8080)

	// Database settings
	config.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/mydb")
	config.DatabaseMaxConns = getEnvAsInt("DATABASE_MAX_CONNS", 10)
	config.DatabaseTimeout = getEnvAsDuration("DATABASE_TIMEOUT", 30*time.Second)

	// Logging settings
	config.LogLevel = getEnv("LOG_LEVEL", "info")
	config.LogFormat = getEnv("LOG_FORMAT", "json")

	// Feature flags
	config.EnableMetrics = getEnvAsBool("ENABLE_METRICS", false)
	config.EnableDebug = getEnvAsBool("ENABLE_DEBUG", false)
	config.EnableCORS = getEnvAsBool("ENABLE_CORS", true)

	// External services
	config.RedisURL = getEnv("REDIS_URL", "redis://localhost:6379")
	config.CacheTimeout = getEnvAsDuration("CACHE_TIMEOUT", 5*time.Minute)

	// Security
	config.JWTSecret = getEnv("JWT_SECRET", "")
	config.SessionSecret = getEnv("SESSION_SECRET", "")

	// Validate required environment variables
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate checks if required environment variables are set
func (c *EnvConfig) Validate() error {
	var errors []string

	// Check required variables based on environment
	if c.AppEnvironment == "production" {
		if c.JWTSecret == "" {
			errors = append(errors, "JWT_SECRET is required in production")
		}
		if c.SessionSecret == "" {
			errors = append(errors, "SESSION_SECRET is required in production")
		}
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLogLevels, c.LogLevel) {
		errors = append(errors, fmt.Sprintf("LOG_LEVEL must be one of: %s", strings.Join(validLogLevels, ", ")))
	}

	// Validate log format
	validLogFormats := []string{"json", "text"}
	if !contains(validLogFormats, c.LogFormat) {
		errors = append(errors, fmt.Sprintf("LOG_FORMAT must be one of: %s", strings.Join(validLogFormats, ", ")))
	}

	// Validate port range
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		errors = append(errors, "SERVER_PORT must be between 1 and 65535")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetServerAddress returns the full server address
func (c *EnvConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort)
}

// IsProduction checks if the application is running in production
func (c *EnvConfig) IsProduction() bool {
	return c.AppEnvironment == "production"
}

// IsDevelopment checks if the application is running in development
func (c *EnvConfig) IsDevelopment() bool {
	return c.AppEnvironment == "development"
}

// IsStaging checks if the application is running in staging
func (c *EnvConfig) IsStaging() bool {
	return c.AppEnvironment == "staging"
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid integer value for %s: %s, using default: %d\n", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid boolean value for %s: %s, using default: %t\n", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid duration value for %s: %s, using default: %v\n", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PrintConfig prints the current configuration (excluding sensitive data)
func (c *EnvConfig) PrintConfig() {
	fmt.Println("=== Configuration ===")
	fmt.Printf("App Name: %s\n", c.AppName)
	fmt.Printf("App Version: %s\n", c.AppVersion)
	fmt.Printf("Environment: %s\n", c.AppEnvironment)
	fmt.Printf("Server: %s\n", c.GetServerAddress())
	fmt.Printf("Database URL: %s\n", maskSensitiveData(c.DatabaseURL))
	fmt.Printf("Database Max Connections: %d\n", c.DatabaseMaxConns)
	fmt.Printf("Database Timeout: %v\n", c.DatabaseTimeout)
	fmt.Printf("Log Level: %s\n", c.LogLevel)
	fmt.Printf("Log Format: %s\n", c.LogFormat)
	fmt.Printf("Metrics Enabled: %t\n", c.EnableMetrics)
	fmt.Printf("Debug Enabled: %t\n", c.EnableDebug)
	fmt.Printf("CORS Enabled: %t\n", c.EnableCORS)
	fmt.Printf("Redis URL: %s\n", maskSensitiveData(c.RedisURL))
	fmt.Printf("Cache Timeout: %v\n", c.CacheTimeout)
	fmt.Printf("JWT Secret Set: %t\n", c.JWTSecret != "")
	fmt.Printf("Session Secret Set: %t\n", c.SessionSecret != "")
}

// maskSensitiveData masks sensitive information in URLs
func maskSensitiveData(url string) string {
	if url == "" {
		return ""
	}

	// Simple masking for demonstration - in production, use proper URL parsing
	if strings.Contains(url, "@") {
		parts := strings.Split(url, "@")
		if len(parts) == 2 {
			return "***@" + parts[1]
		}
	}

	return url
}
