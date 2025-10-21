package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jerrychou/go-practice/concurrency"
)

func main() {
	fmt.Println("Go Concurrency Examples")
	fmt.Println("======================")

	if len(os.Args) > 1 {
		// Run specific example based on command line argument
		example := os.Args[1]
		runSpecificExample(example)
	} else {
		// Run all examples
		runAllExamples()
	}
}

func runSpecificExample(example string) {
	fmt.Printf("Running specific example: %s\n\n", example)

	switch example {
	case "goroutines":
		concurrency.RunAllGoroutineExamples()
	case "channels":
		concurrency.RunAllChannelExamples()
	case "select":
		concurrency.RunAllSelectExamples()
	case "waitgroups":
		concurrency.RunAllWaitGroupExamples()
	case "mutexes":
		concurrency.RunAllMutexExamples()
	case "context":
		concurrency.RunAllContextExamples()
	case "workers":
		concurrency.RunAllWorkerPoolExamples()
	case "fan":
		concurrency.RunAllFanPatternExamples()
	default:
		fmt.Printf("Unknown example: %s\n", example)
		fmt.Println("Available examples: goroutines, channels, select, waitgroups, mutexes, context, workers, fan")
	}
}

func runAllExamples() {
	fmt.Println("Running all concurrency examples...\n")

	// Add delays between examples for better readability
	examples := []struct {
		name     string
		function func()
	}{
		{"Goroutines", concurrency.RunAllGoroutineExamples},
		{"Channels", concurrency.RunAllChannelExamples},
		{"Select Statements", concurrency.RunAllSelectExamples},
		{"WaitGroups", concurrency.RunAllWaitGroupExamples},
		{"Mutexes", concurrency.RunAllMutexExamples},
		{"Context", concurrency.RunAllContextExamples},
		{"Worker Pools", concurrency.RunAllWorkerPoolExamples},
		{"Fan Patterns", concurrency.RunAllFanPatternExamples},
	}

	for i, example := range examples {
		fmt.Printf("=== Example %d: %s ===\n", i+1, example.name)
		example.function()

		if i < len(examples)-1 {
			fmt.Println("\n" + strings.Repeat("=", 50))
			time.Sleep(1 * time.Second) // Pause between examples
		}
	}

	fmt.Println("\n=== All Examples Completed ===")
}

// Interactive mode for running examples
func runInteractiveMode() {
	fmt.Println("Interactive Concurrency Examples")
	fmt.Println("================================")
	fmt.Println("Available examples:")
	fmt.Println("1. Goroutines")
	fmt.Println("2. Channels")
	fmt.Println("3. Select Statements")
	fmt.Println("4. WaitGroups")
	fmt.Println("5. Mutexes")
	fmt.Println("6. Context")
	fmt.Println("7. Worker Pools")
	fmt.Println("8. Fan Patterns")
	fmt.Println("9. Run All")
	fmt.Println("0. Exit")

	for {
		fmt.Print("\nEnter your choice (0-9): ")
		var choice int
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			concurrency.RunAllGoroutineExamples()
		case 2:
			concurrency.RunAllChannelExamples()
		case 3:
			concurrency.RunAllSelectExamples()
		case 4:
			concurrency.RunAllWaitGroupExamples()
		case 5:
			concurrency.RunAllMutexExamples()
		case 6:
			concurrency.RunAllContextExamples()
		case 7:
			concurrency.RunAllWorkerPoolExamples()
		case 8:
			concurrency.RunAllFanPatternExamples()
		case 9:
			runAllExamples()
		case 0:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please enter 0-9.")
		}
	}
}

// Benchmark mode for performance testing
func runBenchmarkMode() {
	fmt.Println("Benchmark Mode")
	fmt.Println("==============")

	// Simple benchmark for different patterns
	patterns := []struct {
		name     string
		function func()
	}{
		{"Goroutines", concurrency.RunAllGoroutineExamples},
		{"Channels", concurrency.RunAllChannelExamples},
		{"Select Statements", concurrency.RunAllSelectExamples},
		{"WaitGroups", concurrency.RunAllWaitGroupExamples},
		{"Mutexes", concurrency.RunAllMutexExamples},
		{"Context", concurrency.RunAllContextExamples},
		{"Worker Pools", concurrency.RunAllWorkerPoolExamples},
		{"Fan Patterns", concurrency.RunAllFanPatternExamples},
	}

	for _, pattern := range patterns {
		start := time.Now()
		pattern.function()
		duration := time.Since(start)
		fmt.Printf("%s completed in: %v\n", pattern.name, duration)
	}
}

// Usage information
func printUsage() {
	fmt.Println("Go Concurrency Examples")
	fmt.Println("=======================")
	fmt.Println("Usage:")
	fmt.Println("  go run concurrency_main.go                    # Run all examples")
	fmt.Println("  go run concurrency_main.go <example>          # Run specific example")
	fmt.Println("  go run concurrency_main.go interactive        # Interactive mode")
	fmt.Println("  go run concurrency_main.go benchmark          # Benchmark mode")
	fmt.Println("")
	fmt.Println("Available examples:")
	fmt.Println("  goroutines  - Basic goroutine spawning and management")
	fmt.Println("  channels    - Channel operations (unbuffered, buffered, directional)")
	fmt.Println("  select      - Select statements for non-blocking operations")
	fmt.Println("  waitgroups  - WaitGroup synchronization patterns")
	fmt.Println("  mutexes     - Mutex and RWMutex for shared state protection")
	fmt.Println("  context     - Context for cancellation and timeouts")
	fmt.Println("  workers     - Worker pool patterns for concurrent processing")
	fmt.Println("  fan         - Fan-in/Fan-out data pipeline patterns")
}

// Check if running in interactive mode
func isInteractiveMode() bool {
	return len(os.Args) > 1 && os.Args[1] == "interactive"
}

// Check if running in benchmark mode
func isBenchmarkMode() bool {
	return len(os.Args) > 1 && os.Args[1] == "benchmark"
}

// Check if help is requested
func isHelpRequested() bool {
	return len(os.Args) > 1 && (os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help")
}

// Initialize the application
func init() {
	// Check for special modes
	if isHelpRequested() {
		printUsage()
		os.Exit(0)
	}

	if isInteractiveMode() {
		runInteractiveMode()
		os.Exit(0)
	}

	if isBenchmarkMode() {
		runBenchmarkMode()
		os.Exit(0)
	}
}

// Example of how to run specific concurrency patterns programmatically
func demonstrateSpecificPatterns() {
	fmt.Println("Demonstrating Specific Patterns")
	fmt.Println("==============================")

	// Example 1: Simple goroutine with channel communication
	fmt.Println("\n1. Simple Goroutine Communication:")
	ch := make(chan string)
	go func() {
		ch <- "Hello from goroutine!"
	}()
	message := <-ch
	fmt.Printf("Received: %s\n", message)

	// Example 2: Worker pool with results
	fmt.Println("\n2. Simple Worker Pool:")
	jobs := make(chan int, 3)
	results := make(chan int, 3)

	// Start worker
	go func() {
		for job := range jobs {
			results <- job * 2
		}
	}()

	// Send jobs
	for i := 1; i <= 3; i++ {
		jobs <- i
	}
	close(jobs)

	// Collect results
	for i := 1; i <= 3; i++ {
		result := <-results
		fmt.Printf("Job %d result: %d\n", i, result)
	}

	// Example 3: Select with timeout
	fmt.Println("\n3. Select with Timeout:")
	timeout := time.After(100 * time.Millisecond)
	select {
	case <-timeout:
		fmt.Println("Timeout occurred")
	default:
		fmt.Println("No timeout")
	}
}

// Run the main function with proper initialization
func runMain() {
	// Check if we should run the demonstration
	if len(os.Args) > 1 && os.Args[1] == "demo" {
		demonstrateSpecificPatterns()
		return
	}

	// Run the main examples
	if len(os.Args) > 1 {
		runSpecificExample(os.Args[1])
	} else {
		runAllExamples()
	}
}
