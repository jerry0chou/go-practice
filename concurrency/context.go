package concurrency

import (
	"context"
	"fmt"
	"time"
)

// BasicContext demonstrates basic context usage
func BasicContext() {
	fmt.Println("=== Basic Context ===")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Start a goroutine that respects context
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context cancelled, stopping goroutine")
				return
			default:
				fmt.Println("Working...")
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// Let it run for a bit
	time.Sleep(300 * time.Millisecond)
	fmt.Println("Main function continues")
}

// ContextWithTimeout demonstrates context with timeout
func ContextWithTimeout() {
	fmt.Println("\n=== Context with Timeout ===")

	// Create context with 1 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start a goroutine that might take longer
	go func() {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Long operation completed")
		case <-ctx.Done():
			fmt.Printf("Operation cancelled: %v\n", ctx.Err())
		}
	}()

	// Wait for context to timeout
	<-ctx.Done()
	fmt.Printf("Context expired: %v\n", ctx.Err())
}

// ContextWithDeadline demonstrates context with deadline
func ContextWithDeadline() {
	fmt.Println("\n=== Context with Deadline ===")

	// Create context with deadline
	deadline := time.Now().Add(500 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Start goroutine that checks deadline
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Deadline reached: %v\n", ctx.Err())
				return
			default:
				fmt.Println("Working until deadline...")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Wait for deadline
	<-ctx.Done()
	fmt.Printf("Context deadline exceeded: %v\n", ctx.Err())
}

// ContextWithCancellation demonstrates manual cancellation
func ContextWithCancellation() {
	fmt.Println("\n=== Context with Manual Cancellation ===")

	ctx, cancel := context.WithCancel(context.Background())

	// Start multiple goroutines
	for i := 1; i <= 3; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Goroutine %d: cancelled\n", id)
					return
				default:
					fmt.Printf("Goroutine %d: working...\n", id)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(i)
	}

	// Cancel after some time
	time.Sleep(300 * time.Millisecond)
	fmt.Println("Cancelling all goroutines...")
	cancel()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("All goroutines should be cancelled")
}

// ContextWithValue demonstrates context with values
func ContextWithValue() {
	fmt.Println("\n=== Context with Values ===")

	// Create context with values
	ctx := context.WithValue(context.Background(), "userID", "12345")
	ctx = context.WithValue(ctx, "requestID", "req-001")

	// Start goroutine that uses context values
	go func() {
		userID := ctx.Value("userID")
		requestID := ctx.Value("requestID")

		fmt.Printf("Processing request %s for user %s\n", requestID, userID)

		// Simulate work
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Request processing completed")
	}()

	time.Sleep(300 * time.Millisecond)
}

// ContextWithNestedCancellation demonstrates nested context cancellation
func ContextWithNestedCancellation() {
	fmt.Println("\n=== Context with Nested Cancellation ===")

	// Create parent context
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Create child context with timeout
	childCtx, childCancel := context.WithTimeout(parentCtx, 500*time.Millisecond)
	defer childCancel()

	// Start goroutine that uses child context
	go func() {
		for {
			select {
			case <-childCtx.Done():
				fmt.Printf("Child context cancelled: %v\n", childCtx.Err())
				return
			default:
				fmt.Println("Child goroutine working...")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Wait for child context to timeout
	<-childCtx.Done()
	fmt.Printf("Child context expired: %v\n", childCtx.Err())
}

// ContextWithTimeoutAndGracefulShutdown demonstrates graceful shutdown
func ContextWithTimeoutAndGracefulShutdown() {
	fmt.Println("\n=== Context with Graceful Shutdown ===")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start goroutine that does cleanup
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Starting graceful shutdown...")
			time.Sleep(200 * time.Millisecond) // Simulate cleanup
			fmt.Println("Graceful shutdown completed")
		}
	}()

	// Wait for timeout
	<-ctx.Done()
	fmt.Printf("Shutdown completed: %v\n", ctx.Err())
}

// ContextWithErrorHandling demonstrates error handling with context
func ContextWithErrorHandling() {
	fmt.Println("\n=== Context with Error Handling ===")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Start goroutine that might return error
	go func() {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Operation would succeed")
		case <-ctx.Done():
			fmt.Printf("Operation cancelled due to: %v\n", ctx.Err())
		}
	}()

	// Wait for context
	<-ctx.Done()

	// Handle different error types
	switch ctx.Err() {
	case context.DeadlineExceeded:
		fmt.Println("Operation timed out")
	case context.Canceled:
		fmt.Println("Operation was cancelled")
	default:
		fmt.Println("Unknown error")
	}
}

// ContextWithMultipleOperations demonstrates context with multiple operations
func ContextWithMultipleOperations() {
	fmt.Println("\n=== Context with Multiple Operations ===")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start multiple operations
	for i := 1; i <= 3; i++ {
		go func(id int) {
			select {
			case <-time.After(time.Duration(id) * 200 * time.Millisecond):
				fmt.Printf("Operation %d completed\n", id)
			case <-ctx.Done():
				fmt.Printf("Operation %d cancelled: %v\n", id, ctx.Err())
			}
		}(i)
	}

	// Wait for context
	<-ctx.Done()
	fmt.Printf("All operations affected by: %v\n", ctx.Err())
}

// ContextWithResourceCleanup demonstrates resource cleanup with context
func ContextWithResourceCleanup() {
	fmt.Println("\n=== Context with Resource Cleanup ===")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Simulate resource acquisition
	resources := make([]string, 0)

	go func() {
		defer func() {
			// Cleanup resources
			fmt.Printf("Cleaning up %d resources\n", len(resources))
			resources = nil
		}()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context cancelled, cleaning up resources")
				return
			default:
				resource := fmt.Sprintf("resource-%d", len(resources)+1)
				resources = append(resources, resource)
				fmt.Printf("Acquired resource: %s\n", resource)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Wait for context
	<-ctx.Done()
	fmt.Printf("Resource cleanup completed: %v\n", ctx.Err())
}

// ContextWithSelect demonstrates context with select statements
func ContextWithSelect() {
	fmt.Println("\n=== Context with Select ===")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Start goroutines that send data
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch1 <- "Data from ch1"
	}()

	go func() {
		time.Sleep(400 * time.Millisecond)
		ch2 <- "Data from ch2"
	}()

	// Use select with context
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received: %s\n", msg2)
		case <-ctx.Done():
			fmt.Printf("Context cancelled: %v\n", ctx.Err())
			return
		}
	}
}

// ContextWithHTTPTimeout demonstrates context for HTTP-like operations
func ContextWithHTTPTimeout() {
	fmt.Println("\n=== Context with HTTP-like Timeout ===")

	// Simulate HTTP request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Simulate HTTP request
	go func() {
		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("HTTP request completed")
		case <-ctx.Done():
			fmt.Printf("HTTP request cancelled: %v\n", ctx.Err())
		}
	}()

	// Wait for timeout
	<-ctx.Done()
	fmt.Printf("HTTP request timed out: %v\n", ctx.Err())
}

// RunAllContextExamples runs all context examples
func RunAllContextExamples() {
	fmt.Println("Running Context Examples...")

	BasicContext()
	ContextWithTimeout()
	ContextWithDeadline()
	ContextWithCancellation()
	ContextWithValue()
	ContextWithNestedCancellation()
	ContextWithTimeoutAndGracefulShutdown()
	ContextWithErrorHandling()
	ContextWithMultipleOperations()
	ContextWithResourceCleanup()
	ContextWithSelect()
	ContextWithHTTPTimeout()

	fmt.Println("\n=== All Context Examples Completed ===")
}
