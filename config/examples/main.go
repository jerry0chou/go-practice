package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jerrychou/go-practice/config"
)

func main() {
	fmt.Println("=== Configuration Management Package Demo ===")
	fmt.Println("This demonstrates all features of the config package:")
	fmt.Println("- Environment Variables (12-factor app)")
	fmt.Println("- Configuration Files (JSON, TOML, YAML)")
	fmt.Println("- Configuration Validation")
	fmt.Println("- Hot Reloading")

	// Example 1: Environment Variables
	fmt.Println("\n1. Environment Variables Configuration")
	exampleEnvironmentConfig()

	// Example 2: Configuration Files
	fmt.Println("\n2. Configuration Files")
	exampleConfigurationFiles()

	// Example 3: Configuration Validation
	fmt.Println("\n3. Configuration Validation")
	exampleValidation()

	// Example 4: Hot Reloading
	fmt.Println("\n4. Hot Reloading")
	exampleHotReload()

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("For more examples, see the README.md files in this directory.")
}

func exampleEnvironmentConfig() {
	fmt.Println("Setting up environment variables...")

	// Set environment variables
	os.Setenv("APP_NAME", "example-app")
	os.Setenv("APP_VERSION", "1.0.0")
	os.Setenv("APP_ENV", "development")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("ENABLE_DEBUG", "true")

	// Load configuration from environment
	envConfig, err := config.LoadFromEnv()
	if err != nil {
		log.Printf("Failed to load environment configuration: %v", err)
		return
	}

	fmt.Printf("✓ Environment configuration loaded\n")
	fmt.Printf("  App: %s v%s (%s)\n", envConfig.AppName, envConfig.AppVersion, envConfig.AppEnvironment)
	fmt.Printf("  Server: %s\n", envConfig.GetServerAddress())
	fmt.Printf("  Debug: %t\n", envConfig.EnableDebug)

	// Validate
	if err := config.ValidateEnvConfig(envConfig); err != nil {
		fmt.Printf("✗ Validation failed: %v\n", err)
	} else {
		fmt.Printf("✓ Configuration is valid\n")
	}
}

func exampleConfigurationFiles() {
	fmt.Println("Loading configuration files...")

	// Test different file formats
	configs := []string{
		"production.toml",
		"development.yaml",
		"simple.json",
	}

	for _, configFile := range configs {
		// Use the file directly since we're running from the examples directory
		configPath := configFile
		// Check if file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("  ⚠️  %s not found, skipping\n", configFile)
			continue
		}

		// Load configuration
		loader := config.NewConfigLoader(configPath)
		cfg, err := loader.Load()
		if err != nil {
			fmt.Printf("  ✗ Failed to load %s: %v\n", configFile, err)
			continue
		}

		fmt.Printf("  ✓ %s loaded (%s format)\n", configFile, loader.GetConfigType())

		// Validate
		if err := config.ValidateFileConfig(cfg); err != nil {
			fmt.Printf("    ✗ Validation failed: %v\n", err)
		} else {
			fmt.Printf("    ✓ Valid configuration\n")
		}
	}
}

func exampleValidation() {
	fmt.Println("Testing configuration validation...")

	// Test valid configuration
	fmt.Println("  Testing valid configuration...")
	validConfig := &config.FileConfig{
		App: config.AppConfig{
			Name:        "test-app",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Database: config.DatabaseConfig{
			URL:            "postgres://localhost:5432/testdb",
			MaxConnections: 10,
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}

	if err := config.ValidateFileConfig(validConfig); err != nil {
		fmt.Printf("    ✗ Valid config failed validation: %v\n", err)
	} else {
		fmt.Printf("    ✓ Valid configuration passed\n")
	}

	// Test invalid configuration
	fmt.Println("  Testing invalid configuration...")
	invalidConfig := &config.FileConfig{
		App: config.AppConfig{
			Name:        "",        // Invalid: empty
			Version:     "invalid", // Invalid: wrong format
			Environment: "invalid", // Invalid: not in enum
		},
		Server: config.ServerConfig{
			Port: 99999, // Invalid: out of range
		},
		Database: config.DatabaseConfig{
			URL:            "invalid-url", // Invalid: not a valid DB URL
			MaxConnections: 0,             // Invalid: must be at least 1
		},
		Logging: config.LoggingConfig{
			Level:  "invalid", // Invalid: not in enum
			Format: "invalid", // Invalid: not in enum
		},
	}

	if err := config.ValidateFileConfig(invalidConfig); err != nil {
		fmt.Printf("    ✓ Invalid config correctly failed: %v\n", err)
	} else {
		fmt.Printf("    ✗ Invalid config should have failed\n")
	}
}

func exampleHotReload() {
	fmt.Println("Setting up hot reload demonstration...")

	// Create temporary directory
	tempDir := "temp_hotreload"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "hotreload.toml")

	// Create initial configuration
	if err := config.CreateDefaultConfig(configPath); err != nil {
		log.Printf("Failed to create test config: %v", err)
		return
	}

	// Create reloadable configuration
	validator := config.CreateDefaultSchema()
	reloadableConfig, err := config.NewReloadableConfig(configPath, &config.FileConfig{}, validator)
	if err != nil {
		log.Printf("Failed to create reloadable config: %v", err)
		return
	}

	// Add callback
	reloadableConfig.AddCallback(func(cfg interface{}) error {
		fmt.Println("    🔄 Configuration reloaded!")
		if fc, ok := cfg.(*config.FileConfig); ok {
			fmt.Printf("       App: %s, Port: %d\n", fc.App.Name, fc.Server.Port)
		}
		return nil
	})

	// Initial load
	if err := reloadableConfig.Reload(); err != nil {
		log.Printf("Failed to load initial config: %v", err)
		return
	}

	fmt.Println("  ✓ Initial configuration loaded")

	// Create hot reload manager
	manager := config.NewHotReloadManager()

	// Add configuration to hot reload
	if err := manager.AddConfig("demo", configPath, reloadableConfig.Reload); err != nil {
		log.Printf("Failed to add config to hot reload: %v", err)
		return
	}

	// Start hot reloading
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := manager.StartAll(ctx); err != nil {
		log.Printf("Failed to start hot reload: %v", err)
		return
	}

	fmt.Println("  ✓ Hot reload started")

	// Simulate configuration change
	time.Sleep(1 * time.Second)

	// Load and modify configuration
	loader := config.NewConfigLoader(configPath)
	currentConfig, err := loader.Load()
	if err != nil {
		log.Printf("Failed to load current config: %v", err)
		return
	}

	// Modify configuration
	currentConfig.App.Name = "hotreload-demo"
	currentConfig.Server.Port = 9999

	// Save modified configuration
	if err := loader.Save(currentConfig); err != nil {
		log.Printf("Failed to save modified config: %v", err)
		return
	}

	fmt.Println("  ✓ Configuration file modified")

	// Wait for reload
	time.Sleep(2 * time.Second)

	// Stop hot reload
	if err := manager.StopAll(); err != nil {
		log.Printf("Failed to stop hot reload: %v", err)
	} else {
		fmt.Println("  ✓ Hot reload stopped")
	}
}
