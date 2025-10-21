package concurrency

import (
	"fmt"
	"sync"
	"time"
)

// BasicWaitGroup demonstrates basic WaitGroup usage
func BasicWaitGroup() {
	fmt.Println("=== Basic WaitGroup ===")

	var wg sync.WaitGroup

	// Start multiple goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment counter
		go func(id int) {
			defer wg.Done() // Decrement counter when done
			fmt.Printf("Goroutine %d starting\n", id)
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All goroutines completed")
}

// WaitGroupWithErrorHandling demonstrates WaitGroup with error handling
func WaitGroupWithErrorHandling() {
	fmt.Println("\n=== WaitGroup with Error Handling ===")

	var wg sync.WaitGroup
	errors := make(chan error, 3)

	// Start goroutines that might fail
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate work that might fail
			if id == 2 {
				errors <- fmt.Errorf("goroutine %d failed", id)
				return
			}

			fmt.Printf("Goroutine %d completed successfully\n", id)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		fmt.Printf("Error: %v\n", err)
	}
}

// WaitGroupWithResults demonstrates collecting results from goroutines
func WaitGroupWithResults() {
	fmt.Println("\n=== WaitGroup with Results ===")

	var wg sync.WaitGroup
	results := make(chan int, 3)

	// Start goroutines that produce results
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate some computation
			result := id * id
			fmt.Printf("Goroutine %d computed: %d\n", id, result)
			results <- result
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Collect and process results
	sum := 0
	for result := range results {
		sum += result
	}
	fmt.Printf("Sum of all results: %d\n", sum)
}

// WaitGroupWithTimeout demonstrates WaitGroup with timeout
func WaitGroupWithTimeout() {
	fmt.Println("\n=== WaitGroup with Timeout ===")

	var wg sync.WaitGroup
	done := make(chan bool)

	// Start goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d starting\n", id)
			time.Sleep(time.Duration(id) * 200 * time.Millisecond)
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	// Start timeout goroutine
	go func() {
		time.Sleep(500 * time.Millisecond)
		done <- true
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		fmt.Println("Timeout reached, some goroutines may still be running")
	case <-func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			wg.Wait()
			close(ch)
		}()
		return ch
	}():
		fmt.Println("All goroutines completed successfully")
	}
}

// WaitGroupWithPanicRecovery demonstrates WaitGroup with panic recovery
func WaitGroupWithPanicRecovery() {
	fmt.Println("\n=== WaitGroup with Panic Recovery ===")

	var wg sync.WaitGroup

	// Start goroutines with panic recovery
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Goroutine %d recovered from panic: %v\n", id, r)
				}
			}()

			fmt.Printf("Goroutine %d starting\n", id)

			// Simulate panic in one goroutine
			if id == 2 {
				panic("Something went wrong!")
			}

			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines completed (with panic recovery)")
}

// WaitGroupWithDefer demonstrates proper defer usage with WaitGroup
func WaitGroupWithDefer() {
	fmt.Println("\n=== WaitGroup with Defer ===")

	var wg sync.WaitGroup

	// Start goroutines with proper defer usage
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done() // Always call Done, even if panic occurs
			defer fmt.Printf("Goroutine %d cleanup completed\n", id)

			fmt.Printf("Goroutine %d starting work\n", id)
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Goroutine %d work completed\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines completed with proper cleanup")
}

// WaitGroupWithNestedGoroutines demonstrates WaitGroup with nested goroutines
func WaitGroupWithNestedGoroutines() {
	fmt.Println("\n=== WaitGroup with Nested Goroutines ===")

	var wg sync.WaitGroup

	// Start main goroutines
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(mainID int) {
			defer wg.Done()

			var nestedWg sync.WaitGroup

			// Start nested goroutines
			for j := 1; j <= 2; j++ {
				nestedWg.Add(1)
				go func(nestedID int) {
					defer nestedWg.Done()
					fmt.Printf("Nested goroutine %d.%d working\n", mainID, nestedID)
					time.Sleep(100 * time.Millisecond)
				}(j)
			}

			// Wait for nested goroutines
			nestedWg.Wait()
			fmt.Printf("Main goroutine %d completed\n", mainID)
		}(i)
	}

	wg.Wait()
	fmt.Println("All main goroutines completed")
}

// WaitGroupWithConditionalExecution demonstrates conditional goroutine execution
func WaitGroupWithConditionalExecution() {
	fmt.Println("\n=== WaitGroup with Conditional Execution ===")

	var wg sync.WaitGroup

	// Start goroutines conditionally
	for i := 1; i <= 5; i++ {
		// Only start goroutines for even numbers
		if i%2 == 0 {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				fmt.Printf("Even goroutine %d working\n", id)
				time.Sleep(100 * time.Millisecond)
			}(i)
		} else {
			fmt.Printf("Skipping odd number %d\n", i)
		}
	}

	wg.Wait()
	fmt.Println("All even goroutines completed")
}

// WaitGroupWithDynamicAddition demonstrates dynamically adding goroutines
func WaitGroupWithDynamicAddition() {
	fmt.Println("\n=== WaitGroup with Dynamic Addition ===")

	var wg sync.WaitGroup

	// Start initial goroutines
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Initial goroutine %d starting\n", id)

			// Dynamically add more goroutines
			if id == 1 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					fmt.Println("Dynamically added goroutine working")
					time.Sleep(100 * time.Millisecond)
				}()
			}

			time.Sleep(200 * time.Millisecond)
			fmt.Printf("Initial goroutine %d finished\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines (including dynamic ones) completed")
}

// WaitGroupWithResourceCleanup demonstrates resource cleanup with WaitGroup
func WaitGroupWithResourceCleanup() {
	fmt.Println("\n=== WaitGroup with Resource Cleanup ===")

	var wg sync.WaitGroup
	resources := make([]string, 0)

	// Start goroutines that use resources
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				fmt.Printf("Cleaning up resources for goroutine %d\n", id)
			}()

			resource := fmt.Sprintf("Resource-%d", id)
			resources = append(resources, resource)
			fmt.Printf("Goroutine %d using resource: %s\n", id, resource)
			time.Sleep(100 * time.Millisecond)
		}(i)
	}

	wg.Wait()
	fmt.Printf("All resources cleaned up: %v\n", resources)
}

// WaitGroupWithProgressTracking demonstrates progress tracking with WaitGroup
func WaitGroupWithProgressTracking() {
	fmt.Println("\n=== WaitGroup with Progress Tracking ===")

	var wg sync.WaitGroup
	progress := make(chan int, 5)

	// Start goroutines that report progress
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 1; j <= 3; j++ {
				time.Sleep(100 * time.Millisecond)
				progress <- id
				fmt.Printf("Goroutine %d progress: %d/3\n", id, j)
			}
		}(i)
	}

	// Start progress tracker
	go func() {
		completed := make(map[int]int)
		for id := range progress {
			completed[id]++
			if completed[id] == 3 {
				fmt.Printf("Goroutine %d completed all tasks\n", id)
			}
		}
	}()

	wg.Wait()
	close(progress)
	fmt.Println("All goroutines completed with progress tracking")
}

// RunAllWaitGroupExamples runs all WaitGroup examples
func RunAllWaitGroupExamples() {
	fmt.Println("Running WaitGroup Examples...")

	BasicWaitGroup()
	WaitGroupWithErrorHandling()
	WaitGroupWithResults()
	WaitGroupWithTimeout()
	WaitGroupWithPanicRecovery()
	WaitGroupWithDefer()
	WaitGroupWithNestedGoroutines()
	WaitGroupWithConditionalExecution()
	WaitGroupWithDynamicAddition()
	WaitGroupWithResourceCleanup()
	WaitGroupWithProgressTracking()

	fmt.Println("\n=== All WaitGroup Examples Completed ===")
}
