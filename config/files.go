package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// FileConfig represents configuration loaded from files
type FileConfig struct {
	App      AppConfig      `json:"app" yaml:"app" toml:"app"`
	Server   ServerConfig   `json:"server" yaml:"server" toml:"server"`
	Database DatabaseConfig `json:"database" yaml:"database" toml:"database"`
	Logging  LoggingConfig  `json:"logging" yaml:"logging" toml:"logging"`
	Features FeatureConfig  `json:"features" yaml:"features" toml:"features"`
	Services ServiceConfig  `json:"services" yaml:"services" toml:"services"`
	Security SecurityConfig `json:"security" yaml:"security" toml:"security"`
}

type AppConfig struct {
	Name        string `json:"name" yaml:"name" toml:"name"`
	Version     string `json:"version" yaml:"version" toml:"version"`
	Environment string `json:"environment" yaml:"environment" toml:"environment"`
	Debug       bool   `json:"debug" yaml:"debug" toml:"debug"`
}

type ServerConfig struct {
	Host         string        `json:"host" yaml:"host" toml:"host"`
	Port         int           `json:"port" yaml:"port" toml:"port"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout"`
}

type DatabaseConfig struct {
	URL               string        `json:"url" yaml:"url" toml:"url"`
	MaxConnections    int           `json:"max_connections" yaml:"max_connections" toml:"max_connections"`
	MinConnections    int           `json:"min_connections" yaml:"min_connections" toml:"min_connections"`
	ConnectionTimeout time.Duration `json:"connection_timeout" yaml:"connection_timeout" toml:"connection_timeout"`
	QueryTimeout      time.Duration `json:"query_timeout" yaml:"query_timeout" toml:"query_timeout"`
	SSLMode           string        `json:"ssl_mode" yaml:"ssl_mode" toml:"ssl_mode"`
}

type LoggingConfig struct {
	Level    string `json:"level" yaml:"level" toml:"level"`
	Format   string `json:"format" yaml:"format" toml:"format"`
	Output   string `json:"output" yaml:"output" toml:"output"`
	Filename string `json:"filename" yaml:"filename" toml:"filename"`
	MaxSize  int    `json:"max_size" yaml:"max_size" toml:"max_size"`
	MaxAge   int    `json:"max_age" yaml:"max_age" toml:"max_age"`
	Compress bool   `json:"compress" yaml:"compress" toml:"compress"`
}

type FeatureConfig struct {
	EnableMetrics   bool `json:"enable_metrics" yaml:"enable_metrics" toml:"enable_metrics"`
	EnableCORS      bool `json:"enable_cors" yaml:"enable_cors" toml:"enable_cors"`
	EnableCache     bool `json:"enable_cache" yaml:"enable_cache" toml:"enable_cache"`
	EnableRateLimit bool `json:"enable_rate_limit" yaml:"enable_rate_limit" toml:"enable_rate_limit"`
}

type ServiceConfig struct {
	Redis RedisConfig `json:"redis" yaml:"redis" toml:"redis"`
	Cache CacheConfig `json:"cache" yaml:"cache" toml:"cache"`
}

type RedisConfig struct {
	URL      string        `json:"url" yaml:"url" toml:"url"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout" toml:"timeout"`
	PoolSize int           `json:"pool_size" yaml:"pool_size" toml:"pool_size"`
}

type CacheConfig struct {
	TTL      time.Duration `json:"ttl" yaml:"ttl" toml:"ttl"`
	MaxSize  int           `json:"max_size" yaml:"max_size" toml:"max_size"`
	Strategy string        `json:"strategy" yaml:"strategy" toml:"strategy"`
}

type SecurityConfig struct {
	JWTSecret     string        `json:"jwt_secret" yaml:"jwt_secret" toml:"jwt_secret"`
	SessionSecret string        `json:"session_secret" yaml:"session_secret" toml:"session_secret"`
	TokenExpiry   time.Duration `json:"token_expiry" yaml:"token_expiry" toml:"token_expiry"`
	BCryptCost    int           `json:"bcrypt_cost" yaml:"bcrypt_cost" toml:"bcrypt_cost"`
}

// ConfigLoader handles loading configuration from various file formats
type ConfigLoader struct {
	configPath string
	configType string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configPath string) *ConfigLoader {
	ext := strings.ToLower(filepath.Ext(configPath))
	configType := "json" // default

	switch ext {
	case ".yaml", ".yml":
		configType = "yaml"
	case ".toml":
		configType = "toml"
	case ".json":
		configType = "json"
	}

	return &ConfigLoader{
		configPath: configPath,
		configType: configType,
	}
}

// Load loads configuration from file
func (cl *ConfigLoader) Load() (*FileConfig, error) {
	if _, err := os.Stat(cl.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", cl.configPath)
	}

	data, err := ioutil.ReadFile(cl.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &FileConfig{}

	switch cl.configType {
	case "yaml":
		err = yaml.Unmarshal(data, config)
	case "toml":
		err = toml.Unmarshal(data, config)
	case "json":
		err = json.Unmarshal(data, config)
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", cl.configType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// Save saves configuration to file
func (cl *ConfigLoader) Save(config *FileConfig) error {
	var data []byte
	var err error

	switch cl.configType {
	case "yaml":
		data, err = yaml.Marshal(config)
	case "toml":
		data, err = toml.Marshal(config)
	case "json":
		data, err = json.MarshalIndent(config, "", "  ")
	default:
		return fmt.Errorf("unsupported config file format: %s", cl.configType)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return ioutil.WriteFile(cl.configPath, data, 0644)
}

// GetConfigType returns the detected configuration file type
func (cl *ConfigLoader) GetConfigType() string {
	return cl.configType
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(configPath string) error {
	ext := strings.ToLower(filepath.Ext(configPath))

	// Determine config type based on file extension
	var configType string
	switch ext {
	case ".yaml", ".yml":
		configType = "yaml"
	case ".toml":
		configType = "toml"
	case ".json":
		configType = "json"
	default:
		configType = "json"
	}

	config := &FileConfig{
		App: AppConfig{
			Name:        "go-practice",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			URL:               "postgres://localhost:5432/mydb",
			MaxConnections:    10,
			MinConnections:    1,
			ConnectionTimeout: 30 * time.Second,
			QueryTimeout:      30 * time.Second,
			SSLMode:           "disable",
		},
		Logging: LoggingConfig{
			Level:    "info",
			Format:   "json",
			Output:   "stdout",
			Filename: "app.log",
			MaxSize:  100,
			MaxAge:   7,
			Compress: true,
		},
		Features: FeatureConfig{
			EnableMetrics:   false,
			EnableCORS:      true,
			EnableCache:     true,
			EnableRateLimit: false,
		},
		Services: ServiceConfig{
			Redis: RedisConfig{
				URL:      "redis://localhost:6379",
				Timeout:  5 * time.Second,
				PoolSize: 10,
			},
			Cache: CacheConfig{
				TTL:      5 * time.Minute,
				MaxSize:  1000,
				Strategy: "lru",
			},
		},
		Security: SecurityConfig{
			JWTSecret:     "",
			SessionSecret: "",
			TokenExpiry:   24 * time.Hour,
			BCryptCost:    12,
		},
	}

	loader := NewConfigLoader(configPath)
	// Use configType to ensure it's not unused
	_ = configType
	return loader.Save(config)
}

// Validate validates the configuration
func (fc *FileConfig) Validate() error {
	var errors []string

	// Validate app configuration
	if fc.App.Name == "" {
		errors = append(errors, "app.name is required")
	}
	if fc.App.Version == "" {
		errors = append(errors, "app.version is required")
	}

	// Validate server configuration
	if fc.Server.Port < 1 || fc.Server.Port > 65535 {
		errors = append(errors, "server.port must be between 1 and 65535")
	}

	// Validate database configuration
	if fc.Database.URL == "" {
		errors = append(errors, "database.url is required")
	}
	if fc.Database.MaxConnections < 1 {
		errors = append(errors, "database.max_connections must be at least 1")
	}

	// Validate logging configuration
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLogLevels, fc.Logging.Level) {
		errors = append(errors, fmt.Sprintf("logging.level must be one of: %s", strings.Join(validLogLevels, ", ")))
	}

	validLogFormats := []string{"json", "text"}
	if !contains(validLogFormats, fc.Logging.Format) {
		errors = append(errors, fmt.Sprintf("logging.format must be one of: %s", strings.Join(validLogFormats, ", ")))
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetServerAddress returns the full server address
func (fc *FileConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", fc.Server.Host, fc.Server.Port)
}

// IsProduction checks if the application is running in production
func (fc *FileConfig) IsProduction() bool {
	return fc.App.Environment == "production"
}

// IsDevelopment checks if the application is running in development
func (fc *FileConfig) IsDevelopment() bool {
	return fc.App.Environment == "development"
}

// PrintConfig prints the current configuration (excluding sensitive data)
func (fc *FileConfig) PrintConfig() {
	fmt.Println("=== File Configuration ===")
	fmt.Printf("App Name: %s\n", fc.App.Name)
	fmt.Printf("App Version: %s\n", fc.App.Version)
	fmt.Printf("Environment: %s\n", fc.App.Environment)
	fmt.Printf("Debug: %t\n", fc.App.Debug)
	fmt.Printf("Server: %s\n", fc.GetServerAddress())
	fmt.Printf("Database URL: %s\n", maskSensitiveData(fc.Database.URL))
	fmt.Printf("Database Max Connections: %d\n", fc.Database.MaxConnections)
	fmt.Printf("Log Level: %s\n", fc.Logging.Level)
	fmt.Printf("Log Format: %s\n", fc.Logging.Format)
	fmt.Printf("Metrics Enabled: %t\n", fc.Features.EnableMetrics)
	fmt.Printf("CORS Enabled: %t\n", fc.Features.EnableCORS)
	fmt.Printf("Cache Enabled: %t\n", fc.Features.EnableCache)
	fmt.Printf("Redis URL: %s\n", maskSensitiveData(fc.Services.Redis.URL))
	fmt.Printf("Cache TTL: %v\n", fc.Services.Cache.TTL)
	fmt.Printf("JWT Secret Set: %t\n", fc.Security.JWTSecret != "")
	fmt.Printf("Session Secret Set: %t\n", fc.Security.SessionSecret != "")
}
