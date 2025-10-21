package concurrency

import (
	"fmt"
	"time"
)

// BasicSelect demonstrates basic select statement usage
func BasicSelect() {
	fmt.Println("=== Basic Select Statement ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Start goroutines to send data
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "Message from channel 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "Message from channel 2"
	}()

	// Select will choose the first available channel
	select {
	case msg1 := <-ch1:
		fmt.Printf("Received from ch1: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("Received from ch2: %s\n", msg2)
	}
}

// SelectWithDefault demonstrates select with default case (non-blocking)
func SelectWithDefault() {
	fmt.Println("\n=== Select with Default (Non-blocking) ===")

	ch := make(chan string)

	// Try to receive from channel with default case
	select {
	case msg := <-ch:
		fmt.Printf("Received: %s\n", msg)
	default:
		fmt.Println("No message available, continuing...")
	}

	// Send a message and try again
	go func() {
		ch <- "Hello!"
	}()

	time.Sleep(50 * time.Millisecond)

	select {
	case msg := <-ch:
		fmt.Printf("Received: %s\n", msg)
	default:
		fmt.Println("No message available")
	}
}

// SelectWithTimeout demonstrates timeout using select
func SelectWithTimeout() {
	fmt.Println("\n=== Select with Timeout ===")

	ch := make(chan string)

	// Start a goroutine that takes time
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Slow operation completed"
	}()

	// Wait for result with timeout
	select {
	case result := <-ch:
		fmt.Printf("Received: %s\n", result)
	case <-time.After(1 * time.Second):
		fmt.Println("Operation timed out!")
	}
}

// SelectWithMultipleCases demonstrates multiple cases in select
func SelectWithMultipleCases() {
	fmt.Println("\n=== Select with Multiple Cases ===")

	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	// Start goroutines with different delays
	go func() {
		time.Sleep(300 * time.Millisecond)
		ch1 <- "From channel 1"
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch2 <- "From channel 2"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch3 <- "From channel 3"
	}()

	// Select will choose the first available
	select {
	case msg1 := <-ch1:
		fmt.Printf("Received: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("Received: %s\n", msg2)
	case msg3 := <-ch3:
		fmt.Printf("Received: %s\n", msg3)
	}
}

// SelectInLoop demonstrates using select in a loop
func SelectInLoop() {
	fmt.Println("\n=== Select in Loop ===")

	ch1 := make(chan string)
	ch2 := make(chan string)
	done := make(chan bool)

	// Start goroutines
	go func() {
		for i := 1; i <= 3; i++ {
			ch1 <- fmt.Sprintf("Message %d from ch1", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for i := 1; i <= 3; i++ {
			ch2 <- fmt.Sprintf("Message %d from ch2", i)
			time.Sleep(150 * time.Millisecond)
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		done <- true
	}()

	// Loop until done signal
	for {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received: %s\n", msg2)
		case <-done:
			fmt.Println("Done signal received, exiting loop")
			return
		}
	}
}

// SelectWithSend demonstrates select with send operations
func SelectWithSend() {
	fmt.Println("\n=== Select with Send Operations ===")

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// Try to send to channels
	select {
	case ch1 <- "Message to ch1":
		fmt.Println("Sent to ch1")
	case ch2 <- "Message to ch2":
		fmt.Println("Sent to ch2")
	default:
		fmt.Println("No channel available for sending")
	}

	// Receive the sent message
	select {
	case msg1 := <-ch1:
		fmt.Printf("Received from ch1: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("Received from ch2: %s\n", msg2)
	}
}

// SelectWithCloseDetection demonstrates detecting channel closure
func SelectWithCloseDetection() {
	fmt.Println("\n=== Select with Close Detection ===")

	ch := make(chan string)

	// Start goroutine to send data and close
	go func() {
		ch <- "First message"
		ch <- "Second message"
		close(ch)
	}()

	// Receive until channel is closed
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				fmt.Println("Channel is closed")
				return
			}
			fmt.Printf("Received: %s\n", msg)
		}
	}
}

// SelectWithTicker demonstrates using ticker with select
func SelectWithTicker() {
	fmt.Println("\n=== Select with Ticker ===")

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	ch := make(chan string)

	// Start goroutine to send data
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "Data received!"
	}()

	// Wait for either ticker or data
	select {
	case <-ticker.C:
		fmt.Println("Ticker fired!")
	case msg := <-ch:
		fmt.Printf("Received data: %s\n", msg)
	}
}

// SelectWithTimer demonstrates using timer with select
func SelectWithTimer() {
	fmt.Println("\n=== Select with Timer ===")

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	ch := make(chan string)

	// Start goroutine to send data
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "Data received!"
	}()

	// Wait for either timer or data
	select {
	case <-timer.C:
		fmt.Println("Timer expired!")
	case msg := <-ch:
		fmt.Printf("Received data: %s\n", msg)
	}
}

// SelectWithDoneChannel demonstrates using done channel for cancellation
func SelectWithDoneChannel() {
	fmt.Println("\n=== Select with Done Channel ===")

	ch := make(chan string)
	done := make(chan bool)

	// Start goroutine to send data
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- fmt.Sprintf("Message %d", i)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Start goroutine to signal done after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		done <- true
	}()

	// Receive until done signal
	for {
		select {
		case msg := <-ch:
			fmt.Printf("Received: %s\n", msg)
		case <-done:
			fmt.Println("Done signal received, stopping...")
			return
		}
	}
}

// SelectWithPriority demonstrates priority-based selection
func SelectWithPriority() {
	fmt.Println("\n=== Select with Priority ===")

	highPriority := make(chan string)
	lowPriority := make(chan string)

	// Start goroutines
	go func() {
		time.Sleep(100 * time.Millisecond)
		highPriority <- "High priority message"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		lowPriority <- "Low priority message"
	}()

	// Check high priority first, then low priority
	select {
	case msg := <-highPriority:
		fmt.Printf("High priority: %s\n", msg)
	default:
		select {
		case msg := <-lowPriority:
			fmt.Printf("Low priority: %s\n", msg)
		default:
			fmt.Println("No messages available")
		}
	}
}

// SelectWithNonBlockingSend demonstrates non-blocking send operations
func SelectWithNonBlockingSend() {
	fmt.Println("\n=== Non-blocking Send ===")

	ch := make(chan string, 1)

	// Try to send without blocking
	select {
	case ch <- "First message":
		fmt.Println("First message sent successfully")
	default:
		fmt.Println("Failed to send first message")
	}

	select {
	case ch <- "Second message":
		fmt.Println("Second message sent successfully")
	default:
		fmt.Println("Failed to send second message (channel full)")
	}

	// Receive messages
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch:
			fmt.Printf("Received: %s\n", msg)
		default:
			fmt.Println("No message to receive")
		}
	}
}

// SelectWithRandomSelection demonstrates random selection when multiple cases are ready
func SelectWithRandomSelection() {
	fmt.Println("\n=== Random Selection ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Start goroutines that send at the same time
	go func() {
		ch1 <- "Message from ch1"
	}()

	go func() {
		ch2 <- "Message from ch2"
	}()

	// Give goroutines time to send
	time.Sleep(10 * time.Millisecond)

	// Select will randomly choose between ready channels
	select {
	case msg1 := <-ch1:
		fmt.Printf("Selected ch1: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("Selected ch2: %s\n", msg2)
	}
}

// RunAllSelectExamples runs all select statement examples
func RunAllSelectExamples() {
	fmt.Println("Running Select Statement Examples...")

	BasicSelect()
	SelectWithDefault()
	SelectWithTimeout()
	SelectWithMultipleCases()
	SelectInLoop()
	SelectWithSend()
	SelectWithCloseDetection()
	SelectWithTicker()
	SelectWithTimer()
	SelectWithDoneChannel()
	SelectWithPriority()
	SelectWithNonBlockingSend()
	SelectWithRandomSelection()

	fmt.Println("\n=== All Select Statement Examples Completed ===")
}
