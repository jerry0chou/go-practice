package concurrency

import (
	"fmt"
	"time"
)

// UnbufferedChannels demonstrates unbuffered channels (synchronous communication)
func UnbufferedChannels() {
	fmt.Println("=== Unbuffered Channels ===")

	// Create an unbuffered channel
	ch := make(chan string)

	// Start a goroutine to send data
	go func() {
		fmt.Println("Sending data to channel...")
		ch <- "Hello from goroutine!"
		fmt.Println("Data sent successfully")
	}()

	// Receive data from channel
	fmt.Println("Waiting to receive data...")
	message := <-ch
	fmt.Printf("Received: %s\n", message)
}

// BufferedChannels demonstrates buffered channels (asynchronous communication)
func BufferedChannels() {
	fmt.Println("\n=== Buffered Channels ===")

	// Create a buffered channel with capacity 3
	ch := make(chan int, 3)

	// Send multiple values without blocking (until buffer is full)
	fmt.Println("Sending values to buffered channel...")
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("All values sent to buffer")

	// Receive values
	fmt.Println("Receiving values:")
	for i := 0; i < 3; i++ {
		value := <-ch
		fmt.Printf("Received: %d\n", value)
	}
}

// SendOnlyChannel demonstrates send-only channels
func SendOnlyChannel() {
	fmt.Println("\n=== Send-Only Channel ===")

	ch := make(chan string)

	// Function that only sends to channel
	sendData := func(ch chan<- string) {
		ch <- "Data from send-only function"
		close(ch)
	}

	// Start goroutine with send-only channel
	go sendData(ch)

	// Receive data
	message := <-ch
	fmt.Printf("Received: %s\n", message)
}

// ReceiveOnlyChannel demonstrates receive-only channels
func ReceiveOnlyChannel() {
	fmt.Println("\n=== Receive-Only Channel ===")

	ch := make(chan string)

	// Function that only receives from channel
	receiveData := func(ch <-chan string) {
		for message := range ch {
			fmt.Printf("Received: %s\n", message)
		}
	}

	// Start goroutine with receive-only channel
	go receiveData(ch)

	// Send data
	ch <- "Message 1"
	ch <- "Message 2"
	close(ch)

	// Give time for processing
	time.Sleep(100 * time.Millisecond)
}

// ChannelDirectionExample demonstrates channel direction in function parameters
func ChannelDirectionExample() {
	fmt.Println("\n=== Channel Direction Example ===")

	ch := make(chan int)

	// Producer function (send-only)
	producer := func(ch chan<- int) {
		for i := 1; i <= 5; i++ {
			ch <- i
			fmt.Printf("Sent: %d\n", i)
		}
		close(ch)
	}

	// Consumer function (receive-only)
	consumer := func(ch <-chan int) {
		for value := range ch {
			fmt.Printf("Received: %d\n", value)
		}
	}

	// Start producer and consumer
	go producer(ch)
	consumer(ch)
}

// ChannelWithSelect demonstrates using channels with select statements
func ChannelWithSelect() {
	fmt.Println("\n=== Channel with Select ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Start goroutines to send data
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "From channel 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "From channel 2"
	}()

	// Use select to receive from whichever channel is ready
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received: %s\n", msg2)
		}
	}

}

// ChannelTimeout demonstrates timeout with channels
func ChannelTimeout() {
	fmt.Println("\n=== Channel Timeout ===")

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

// ChannelCloseDetection demonstrates how to detect channel closure
func ChannelCloseDetection() {
	fmt.Println("\n=== Channel Close Detection ===")

	ch := make(chan int)

	// Start goroutine to send data and close channel
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)
	}()

	// Receive data until channel is closed
	for {
		value, ok := <-ch
		if !ok {
			fmt.Println("Channel is closed")
			break
		}
		fmt.Printf("Received: %d\n", value)
	}
}

// ChannelRange demonstrates using range with channels
func ChannelRange() {
	fmt.Println("\n=== Channel Range ===")

	ch := make(chan string)

	// Start goroutine to send data
	go func() {
		ch <- "First message"
		ch <- "Second message"
		ch <- "Third message"
		close(ch)
	}()

	// Use range to receive all values
	for message := range ch {
		fmt.Printf("Received: %s\n", message)
	}
}

// ChannelCapacityExample demonstrates checking channel capacity
func ChannelCapacityExample() {
	fmt.Println("\n=== Channel Capacity Example ===")

	// Unbuffered channel
	unbuffered := make(chan int)
	fmt.Printf("Unbuffered channel capacity: %d\n", cap(unbuffered))

	// Buffered channel
	buffered := make(chan int, 5)
	fmt.Printf("Buffered channel capacity: %d\n", cap(buffered))

	// Check length of buffered channel
	buffered <- 1
	buffered <- 2
	fmt.Printf("Buffered channel length: %d\n", len(buffered))

	// Clean up
	<-buffered
	<-buffered
}

// ChannelWithMultipleReceivers demonstrates multiple receivers on one channel
func ChannelWithMultipleReceivers() {
	fmt.Println("\n=== Multiple Receivers ===")

	ch := make(chan string)

	// Start multiple receivers
	for i := 1; i <= 3; i++ {
		go func(id int) {
			message := <-ch
			fmt.Printf("Receiver %d got: %s\n", id, message)
		}(i)
	}

	// Send one message (only one receiver will get it)
	ch <- "Broadcast message"

	time.Sleep(100 * time.Millisecond)
}

// ChannelWithMultipleSenders demonstrates multiple senders on one channel
func ChannelWithMultipleSenders() {
	fmt.Println("\n=== Multiple Senders ===")

	ch := make(chan string)

	// Start multiple senders
	for i := 1; i <= 3; i++ {
		go func(id int) {
			ch <- fmt.Sprintf("Message from sender %d", id)
		}(i)
	}

	// Receive messages
	for i := 0; i < 3; i++ {
		message := <-ch
		fmt.Printf("Received: %s\n", message)
	}
}

// ChannelWithDone demonstrates using a done channel for cancellation
func ChannelWithDone() {
	fmt.Println("\n=== Channel with Done ===")

	ch := make(chan int)
	done := make(chan bool)

	// Start a goroutine that sends data
	go func() {
		for i := 1; i <= 10; i++ {
			select {
			case ch <- i:
				fmt.Printf("Sent: %d\n", i)
			case <-done:
				fmt.Println("Sender stopped by done signal")
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Start a goroutine that stops after receiving 5 values
	go func() {
		count := 0
		for range ch {
			count++
			if count >= 5 {
				fmt.Println("Stopping after 5 values")
				done <- true
				return
			}
		}
	}()

	// Wait for completion
	time.Sleep(2 * time.Second)
}

// RunAllChannelExamples runs all channel examples
func RunAllChannelExamples() {
	fmt.Println("Running Channel Examples...")

	UnbufferedChannels()
	BufferedChannels()
	SendOnlyChannel()
	ReceiveOnlyChannel()
	ChannelDirectionExample()
	ChannelWithSelect()
	ChannelTimeout()
	ChannelCloseDetection()
	ChannelRange()
	ChannelCapacityExample()
	ChannelWithMultipleReceivers()
	ChannelWithMultipleSenders()
	ChannelWithDone()

	fmt.Println("\n=== All Channel Examples Completed ===")
}
