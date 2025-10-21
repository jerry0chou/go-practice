package concurrency

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// BasicGoroutine demonstrates the simplest way to start a goroutine
func BasicGoroutine() {
	fmt.Println("=== Basic Goroutine Example ===")

	// Start a goroutine with an anonymous function
	go func() {
		fmt.Println("Hello from goroutine!")
	}()

	// Give the goroutine time to execute
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main function continues...")
}

// GoroutineWithFunction demonstrates starting a goroutine with a named function
func GoroutineWithFunction() {
	fmt.Println("\n=== Goroutine with Named Function ===")

	go printNumbers(5)
	time.Sleep(200 * time.Millisecond)
}

// printNumbers is a helper function for goroutines
func printNumbers(n int) {
	for i := 1; i <= n; i++ {
		fmt.Printf("Number: %d\n", i)
		time.Sleep(50 * time.Millisecond)
	}
}

// MultipleGoroutines demonstrates running multiple goroutines concurrently
func MultipleGoroutines() {
	fmt.Println("\n=== Multiple Goroutines ===")

	var wg sync.WaitGroup

	// Start multiple goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d starting\n", id)
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All goroutines completed")
}

// GoroutineLifecycle demonstrates goroutine lifecycle management
func GoroutineLifecycle() {
	fmt.Println("\n=== Goroutine Lifecycle ===")

	// Channel to signal goroutine completion
	done := make(chan bool)

	go func() {
		fmt.Println("Goroutine is working...")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Goroutine finished work")
		done <- true
	}()

	// Wait for goroutine to complete
	<-done
	fmt.Println("Main function received completion signal")
}

// GoroutineWithReturnValue demonstrates how to get return values from goroutines
func GoroutineWithReturnValue() {
	fmt.Println("\n=== Goroutine with Return Value ===")

	// Channel to receive the result
	result := make(chan int)

	go func() {
		// Simulate some computation
		time.Sleep(200 * time.Millisecond)
		result <- 42
	}()

	// Receive the result
	value := <-result
	fmt.Printf("Received value: %d\n", value)
}

// GoroutineWithErrorHandling demonstrates error handling in goroutines
func GoroutineWithErrorHandling() {
	fmt.Println("\n=== Goroutine with Error Handling ===")

	type Result struct {
		Value int
		Error error
	}

	result := make(chan Result)

	go func() {
		// Simulate a computation that might fail
		time.Sleep(100 * time.Millisecond)

		// Simulate success (you could add error conditions here)
		result <- Result{Value: 100, Error: nil}
	}()

	// Handle the result
	res := <-result
	if res.Error != nil {
		fmt.Printf("Error: %v\n", res.Error)
	} else {
		fmt.Printf("Success: %d\n", res.Value)
	}
}

// GoroutineWithTimeout demonstrates timeout handling for goroutines
func GoroutineWithTimeout() {
	fmt.Println("\n=== Goroutine with Timeout ===")

	result := make(chan string)

	// Start a goroutine that might take too long
	go func() {
		time.Sleep(2 * time.Second)
		result <- "Task completed"
	}()

	// Wait for result with timeout
	select {
	case res := <-result:
		fmt.Printf("Received: %s\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Operation timed out!")
	}
}

// GoroutineStats demonstrates how to get information about running goroutines
func GoroutineStats() {
	fmt.Println("\n=== Goroutine Statistics ===")

	// Get current number of goroutines
	numGoroutines := runtime.NumGoroutine()
	fmt.Printf("Current number of goroutines: %d\n", numGoroutines)

	// Start some goroutines
	for i := 0; i < 5; i++ {
		go func(id int) {
			time.Sleep(1 * time.Second)
		}(i)
	}

	// Check goroutine count again
	time.Sleep(100 * time.Millisecond)
	numGoroutines = runtime.NumGoroutine()
	fmt.Printf("Number of goroutines after starting 5: %d\n", numGoroutines)

	// Wait for goroutines to complete
	time.Sleep(2 * time.Second)
	numGoroutines = runtime.NumGoroutine()
	fmt.Printf("Number of goroutines after completion: %d\n", numGoroutines)
}

// GoroutineWithPanicRecovery demonstrates panic recovery in goroutines
func GoroutineWithPanicRecovery() {
	fmt.Println("\n=== Goroutine with Panic Recovery ===")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
		}()

		fmt.Println("Goroutine starting...")
		time.Sleep(100 * time.Millisecond)

		// This will cause a panic
		panic("Something went wrong!")
	}()

	wg.Wait()
	fmt.Println("Main function continues after panic recovery")
}

// GoroutineWithDefer demonstrates defer statements in goroutines
func GoroutineWithDefer() {
	fmt.Println("\n=== Goroutine with Defer ===")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("Goroutine cleanup completed")
		defer fmt.Println("Goroutine finishing...")

		fmt.Println("Goroutine starting work...")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Goroutine work completed")
	}()

	wg.Wait()
	fmt.Println("Main function continues")
}

// RunAllGoroutineExamples runs all goroutine examples
func RunAllGoroutineExamples() {
	fmt.Println("Running Goroutine Examples...")

	BasicGoroutine()
	GoroutineWithFunction()
	MultipleGoroutines()
	GoroutineLifecycle()
	GoroutineWithReturnValue()
	GoroutineWithErrorHandling()
	GoroutineWithTimeout()
	GoroutineStats()
	GoroutineWithPanicRecovery()
	GoroutineWithDefer()

	fmt.Println("\n=== All Goroutine Examples Completed ===")
}
