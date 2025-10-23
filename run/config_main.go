package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("=== Configuration Management Demo ===")
	fmt.Println("This will run the configuration management examples...")

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Path to the config examples
	configExamplesPath := filepath.Join(currentDir, "config", "examples")
	fmt.Println("Configuration Examples Path:", configExamplesPath)
	// Check if the examples directory exists
	if _, err := os.Stat(configExamplesPath); os.IsNotExist(err) {
		fmt.Printf("Config examples directory not found: %s\n", configExamplesPath)
		return
	}

	// Change to the examples directory and run main.go
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = configExamplesPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running configuration examples from: %s\n", configExamplesPath)
	fmt.Println(strings.Repeat("=", 60))

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running config examples: %v\n", err)
		return
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Configuration management demo completed!")
	fmt.Println("\nTo run examples directly:")
	fmt.Printf("  cd %s\n", configExamplesPath)
	fmt.Println("  go run main.go")
}
