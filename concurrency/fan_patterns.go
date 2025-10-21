package concurrency

import (
	"fmt"
	"sync"
	"time"
)

// FanOutBasic demonstrates basic fan-out pattern
func FanOutBasic() {
	fmt.Println("=== Basic Fan-Out Pattern ===")

	// Input channel
	input := make(chan int, 10)

	// Output channels
	output1 := make(chan int, 5)
	output2 := make(chan int, 5)
	output3 := make(chan int, 5)

	// Start fan-out goroutines
	go func() {
		for data := range input {
			fmt.Printf("Fan-out 1: processing %d\n", data)
			output1 <- data * 2
		}
		close(output1)
	}()

	go func() {
		for data := range input {
			fmt.Printf("Fan-out 2: processing %d\n", data)
			output2 <- data * 3
		}
		close(output2)
	}()

	go func() {
		for data := range input {
			fmt.Printf("Fan-out 3: processing %d\n", data)
			output3 <- data * 4
		}
		close(output3)
	}()

	// Send data to input
	go func() {
		for i := 1; i <= 6; i++ {
			input <- i
		}
		close(input)
	}()

	// Collect results from all outputs
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for result := range output1 {
			fmt.Printf("Output 1 result: %d\n", result)
		}
	}()

	go func() {
		defer wg.Done()
		for result := range output2 {
			fmt.Printf("Output 2 result: %d\n", result)
		}
	}()

	go func() {
		defer wg.Done()
		for result := range output3 {
			fmt.Printf("Output 3 result: %d\n", result)
		}
	}()

	wg.Wait()
}

// FanInBasic demonstrates basic fan-in pattern
func FanInBasic() {
	fmt.Println("\n=== Basic Fan-In Pattern ===")

	// Input channels
	input1 := make(chan int, 3)
	input2 := make(chan int, 3)
	input3 := make(chan int, 3)

	// Output channel
	output := make(chan int, 10)

	// Start fan-in goroutine
	go func() {
		var wg sync.WaitGroup
		wg.Add(3)

		// Fan-in from input1
		go func() {
			defer wg.Done()
			for data := range input1 {
				fmt.Printf("Fan-in: received %d from input1\n", data)
				output <- data
			}
		}()

		// Fan-in from input2
		go func() {
			defer wg.Done()
			for data := range input2 {
				fmt.Printf("Fan-in: received %d from input2\n", data)
				output <- data
			}
		}()

		// Fan-in from input3
		go func() {
			defer wg.Done()
			for data := range input3 {
				fmt.Printf("Fan-in: received %d from input3\n", data)
				output <- data
			}
		}()

		wg.Wait()
		close(output)
	}()

	// Send data to inputs
	go func() {
		for i := 1; i <= 3; i++ {
			input1 <- i
		}
		close(input1)
	}()

	go func() {
		for i := 4; i <= 6; i++ {
			input2 <- i
		}
		close(input2)
	}()

	go func() {
		for i := 7; i <= 9; i++ {
			input3 <- i
		}
		close(input3)
	}()

	// Collect results
	for result := range output {
		fmt.Printf("Final result: %d\n", result)
	}
}

// FanOutFanInPipeline demonstrates a complete fan-out/fan-in pipeline
func FanOutFanInPipeline() {
	fmt.Println("\n=== Fan-Out/Fan-In Pipeline ===")

	// Stage 1: Generate data
	generate := func() <-chan int {
		out := make(chan int)
		go func() {
			for i := 1; i <= 10; i++ {
				fmt.Printf("Generator: producing %d\n", i)
				out <- i
			}
			close(out)
		}()
		return out
	}

	// Stage 2: Fan-out processing
	process := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			for data := range in {
				fmt.Printf("Processor: processing %d\n", data)
				time.Sleep(50 * time.Millisecond) // Simulate work
				out <- data * 2
			}
			close(out)
		}()
		return out
	}

	// Stage 3: Fan-in aggregation
	aggregate := func(inputs ...<-chan int) <-chan int {
		out := make(chan int)
		var wg sync.WaitGroup

		for _, input := range inputs {
			wg.Add(1)
			go func(ch <-chan int) {
				defer wg.Done()
				for data := range ch {
					fmt.Printf("Aggregator: received %d\n", data)
					out <- data
				}
			}(input)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}

	// Create pipeline
	gen := generate()

	// Fan-out to multiple processors
	proc1 := process(gen)
	proc2 := process(gen)
	proc3 := process(gen)

	// Fan-in results
	results := aggregate(proc1, proc2, proc3)

	// Collect final results
	sum := 0
	for result := range results {
		sum += result
		fmt.Printf("Final result: %d\n", result)
	}
	fmt.Printf("Total sum: %d\n", sum)
}

// FanOutWithWorkerPool demonstrates fan-out with worker pool
func FanOutWithWorkerPool() {
	fmt.Println("\n=== Fan-Out with Worker Pool ===")

	const numWorkers = 3
	const numJobs = 9

	// Input channel
	input := make(chan int, numJobs)

	// Output channels for each worker
	outputs := make([]chan int, numWorkers)
	for i := range outputs {
		outputs[i] = make(chan int, 3)
	}

	// Start workers
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			for data := range input {
				fmt.Printf("Worker %d: processing %d\n", workerID, data)
				time.Sleep(100 * time.Millisecond)
				outputs[workerID] <- data * (workerID + 1)
			}
			close(outputs[workerID])
		}(w)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			input <- j
		}
		close(input)
	}()

	// Collect results from all workers
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i, output := range outputs {
		go func(workerID int, ch <-chan int) {
			defer wg.Done()
			for result := range ch {
				fmt.Printf("Worker %d result: %d\n", workerID, result)
			}
		}(i, output)
	}

	wg.Wait()
}

// FanInWithSelect demonstrates fan-in using select statements
func FanInWithSelect() {
	fmt.Println("\n=== Fan-In with Select ===")

	// Input channels
	input1 := make(chan string, 2)
	input2 := make(chan string, 2)
	input3 := make(chan string, 2)

	// Output channel
	output := make(chan string, 10)

	// Start fan-in with select
	go func() {
		defer close(output)

		for {
			select {
			case data, ok := <-input1:
				if !ok {
					input1 = nil
				} else {
					fmt.Printf("Fan-in: received %s from input1\n", data)
					output <- data
				}
			case data, ok := <-input2:
				if !ok {
					input2 = nil
				} else {
					fmt.Printf("Fan-in: received %s from input2\n", data)
					output <- data
				}
			case data, ok := <-input3:
				if !ok {
					input3 = nil
				} else {
					fmt.Printf("Fan-in: received %s from input3\n", data)
					output <- data
				}
			}

			// Exit when all inputs are closed
			if input1 == nil && input2 == nil && input3 == nil {
				return
			}
		}
	}()

	// Send data to inputs
	go func() {
		input1 <- "A1"
		input1 <- "A2"
		close(input1)
	}()

	go func() {
		input2 <- "B1"
		input2 <- "B2"
		close(input2)
	}()

	go func() {
		input3 <- "C1"
		input3 <- "C2"
		close(input3)
	}()

	// Collect results
	for result := range output {
		fmt.Printf("Final result: %s\n", result)
	}
}

// FanOutWithErrorHandling demonstrates fan-out with error handling
func FanOutWithErrorHandling() {
	fmt.Println("\n=== Fan-Out with Error Handling ===")

	type Result struct {
		Data  int
		Error error
	}

	input := make(chan int, 5)
	output := make(chan Result, 10)

	// Start fan-out with error handling
	go func() {
		for data := range input {
			fmt.Printf("Processing %d\n", data)

			// Simulate work that might fail
			time.Sleep(50 * time.Millisecond)

			if data%3 == 0 {
				output <- Result{Data: data, Error: fmt.Errorf("processing failed for %d", data)}
			} else {
				output <- Result{Data: data * 2, Error: nil}
			}
		}
		close(output)
	}()

	// Send data
	go func() {
		for i := 1; i <= 6; i++ {
			input <- i
		}
		close(input)
	}()

	// Collect results
	successCount := 0
	errorCount := 0

	for result := range output {
		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
			errorCount++
		} else {
			fmt.Printf("Success: %d\n", result.Data)
			successCount++
		}
	}

	fmt.Printf("Results: %d success, %d errors\n", successCount, errorCount)
}

// FanInWithBuffering demonstrates fan-in with buffering
func FanInWithBuffering() {
	fmt.Println("\n=== Fan-In with Buffering ===")

	// Input channels with different buffer sizes
	input1 := make(chan int, 2)
	input2 := make(chan int, 3)
	input3 := make(chan int, 1)

	// Output channel with larger buffer
	output := make(chan int, 10)

	// Start fan-in
	go func() {
		defer close(output)

		for {
			select {
			case data, ok := <-input1:
				if !ok {
					input1 = nil
				} else {
					output <- data
				}
			case data, ok := <-input2:
				if !ok {
					input2 = nil
				} else {
					output <- data
				}
			case data, ok := <-input3:
				if !ok {
					input3 = nil
				} else {
					output <- data
				}
			}

			if input1 == nil && input2 == nil && input3 == nil {
				return
			}
		}
	}()

	// Send data at different rates
	go func() {
		for i := 1; i <= 3; i++ {
			input1 <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(input1)
	}()

	go func() {
		for i := 4; i <= 6; i++ {
			input2 <- i
			time.Sleep(150 * time.Millisecond)
		}
		close(input2)
	}()

	go func() {
		for i := 7; i <= 8; i++ {
			input3 <- i
			time.Sleep(200 * time.Millisecond)
		}
		close(input3)
	}()

	// Collect results
	for result := range output {
		fmt.Printf("Buffered result: %d\n", result)
	}
}

// FanOutWithLoadBalancing demonstrates fan-out with load balancing
func FanOutWithLoadBalancing() {
	fmt.Println("\n=== Fan-Out with Load Balancing ===")

	const numWorkers = 3
	const numJobs = 12

	input := make(chan int, numJobs)
	outputs := make([]chan int, numWorkers)

	// Initialize output channels
	for i := range outputs {
		outputs[i] = make(chan int, 4)
	}

	// Start workers with different processing times
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			processingTime := time.Duration(workerID+1) * 50 * time.Millisecond

			for data := range input {
				fmt.Printf("Worker %d: processing %d (time: %v)\n", workerID, data, processingTime)
				time.Sleep(processingTime)
				outputs[workerID] <- data * (workerID + 1)
			}
			close(outputs[workerID])
		}(w)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			input <- j
		}
		close(input)
	}()

	// Collect results
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i, output := range outputs {
		go func(workerID int, ch <-chan int) {
			defer wg.Done()
			for result := range ch {
				fmt.Printf("Worker %d result: %d\n", workerID, result)
			}
		}(i, output)
	}

	wg.Wait()
}

// FanInWithPriority demonstrates fan-in with priority handling
func FanInWithPriority() {
	fmt.Println("\n=== Fan-In with Priority ===")

	// High priority input
	highPriority := make(chan string, 3)
	// Low priority input
	lowPriority := make(chan string, 3)

	output := make(chan string, 10)

	// Start fan-in with priority
	go func() {
		defer close(output)

		for {
			// Check high priority first
			select {
			case data, ok := <-highPriority:
				if !ok {
					highPriority = nil
				} else {
					output <- data
					continue
				}
			default:
			}

			// If no high priority, check low priority
			select {
			case data, ok := <-highPriority:
				if !ok {
					highPriority = nil
				} else {
					output <- data
				}
			case data, ok := <-lowPriority:
				if !ok {
					lowPriority = nil
				} else {
					output <- data
				}
			}

			if highPriority == nil && lowPriority == nil {
				return
			}
		}
	}()

	// Send high priority data
	go func() {
		highPriority <- "HIGH-1"
		highPriority <- "HIGH-2"
		close(highPriority)
	}()

	// Send low priority data
	go func() {
		lowPriority <- "LOW-1"
		lowPriority <- "LOW-2"
		close(lowPriority)
	}()

	// Collect results
	for result := range output {
		fmt.Printf("Priority result: %s\n", result)
	}
}

// RunAllFanPatternExamples runs all fan pattern examples
func RunAllFanPatternExamples() {
	fmt.Println("Running Fan Pattern Examples...")

	FanOutBasic()
	//FanInBasic()
	//FanOutFanInPipeline()
	//FanOutWithWorkerPool()
	//FanInWithSelect()
	//FanOutWithErrorHandling()
	//FanInWithBuffering()
	//FanOutWithLoadBalancing()
	//FanInWithPriority()

	fmt.Println("\n=== All Fan Pattern Examples Completed ===")
}
