package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// BasicWorkerPool demonstrates a basic worker pool pattern
func BasicWorkerPool() {
	fmt.Println("=== Basic Worker Pool ===")

	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", id, job)
				time.Sleep(100 * time.Millisecond) // Simulate work
				results <- job * 2
			}
		}(w)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for r := 1; r <= numJobs; r++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}

// WorkerPoolWithWaitGroup demonstrates worker pool with WaitGroup
func WorkerPoolWithWaitGroup() {
	fmt.Println("\n=== Worker Pool with WaitGroup ===")

	const numWorkers = 4
	const numJobs = 8

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", id, job)
				time.Sleep(50 * time.Millisecond)
				results <- job * job
			}
		}(w)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	// Wait for workers to complete
	wg.Wait()
	close(results)

	// Collect results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

// WorkerPoolWithContext demonstrates worker pool with context cancellation
func WorkerPoolWithContext() {
	fmt.Println("\n=== Worker Pool with Context ===")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	const numWorkers = 3
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}
					fmt.Printf("Worker %d processing job %d\n", id, job)
					time.Sleep(100 * time.Millisecond)
					results <- job * 3
				case <-ctx.Done():
					fmt.Printf("Worker %d cancelled\n", id)
					return
				}
			}
		}(w)
	}

	// Send jobs
	go func() {
		for j := 1; j <= 10; j++ {
			select {
			case jobs <- j:
				fmt.Printf("Sent job %d\n", j)
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Wait for context or collect results
	time.Sleep(600 * time.Millisecond)
	close(results)

	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

// WorkerPoolWithErrorHandling demonstrates worker pool with error handling
func WorkerPoolWithErrorHandling() {
	fmt.Println("\n=== Worker Pool with Error Handling ===")

	const numWorkers = 3
	const numJobs = 6

	type Job struct {
		ID   int
		Data string
	}

	type Result struct {
		JobID int
		Value string
		Error error
	}

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", id, job.ID)

				// Simulate work that might fail
				time.Sleep(100 * time.Millisecond)

				var result Result
				if job.ID%3 == 0 {
					result = Result{JobID: job.ID, Error: fmt.Errorf("job %d failed", job.ID)}
				} else {
					result = Result{JobID: job.ID, Value: fmt.Sprintf("processed-%s", job.Data)}
				}

				results <- result
			}
		}(w)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Data: fmt.Sprintf("data-%d", j)}
	}
	close(jobs)

	// Collect results
	for r := 1; r <= numJobs; r++ {
		result := <-results
		if result.Error != nil {
			fmt.Printf("Job %d failed: %v\n", result.JobID, result.Error)
		} else {
			fmt.Printf("Job %d result: %s\n", result.JobID, result.Value)
		}
	}
}

// WorkerPoolWithRateLimiting demonstrates worker pool with rate limiting
func WorkerPoolWithRateLimiting() {
	fmt.Println("\n=== Worker Pool with Rate Limiting ===")

	const numWorkers = 2
	const numJobs = 8

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Rate limiter: process at most 1 job per 200ms per worker
	rateLimiter := time.NewTicker(200 * time.Millisecond)
	defer rateLimiter.Stop()

	// Start workers with rate limiting
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for job := range jobs {
				<-rateLimiter.C // Wait for rate limit
				fmt.Printf("Worker %d processing job %d (rate limited)\n", id, job)
				time.Sleep(50 * time.Millisecond)
				results <- job * 4
			}
		}(w)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for r := 1; r <= numJobs; r++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}

// WorkerPoolWithPriority demonstrates worker pool with priority queues
func WorkerPoolWithPriority() {
	fmt.Println("\n=== Worker Pool with Priority ===")

	type PriorityJob struct {
		ID       int
		Priority int
		Data     string
	}

	const numWorkers = 2
	highPriorityJobs := make(chan PriorityJob, 5)
	lowPriorityJobs := make(chan PriorityJob, 5)
	results := make(chan string, 10)

	// Start workers that prioritize high priority jobs
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for {
				select {
				case job, ok := <-highPriorityJobs:
					if !ok {
						highPriorityJobs = nil
					} else {
						fmt.Printf("Worker %d processing HIGH priority job %d\n", id, job.ID)
						time.Sleep(100 * time.Millisecond)
						results <- fmt.Sprintf("HIGH-%d: %s", job.ID, job.Data)
					}
				default:
					select {
					case job, ok := <-lowPriorityJobs:
						if !ok {
							lowPriorityJobs = nil
						} else {
							fmt.Printf("Worker %d processing LOW priority job %d\n", id, job.ID)
							time.Sleep(100 * time.Millisecond)
							results <- fmt.Sprintf("LOW-%d: %s", job.ID, job.Data)
						}
					case <-time.After(50 * time.Millisecond):
						// No jobs available
					}
				}

				// Exit when both channels are closed
				if highPriorityJobs == nil && lowPriorityJobs == nil {
					return
				}
			}
		}(w)
	}

	// Send jobs with different priorities
	go func() {
		for i := 1; i <= 3; i++ {
			highPriorityJobs <- PriorityJob{ID: i, Priority: 1, Data: fmt.Sprintf("urgent-%d", i)}
		}
		close(highPriorityJobs)
	}()

	go func() {
		for i := 4; i <= 6; i++ {
			lowPriorityJobs <- PriorityJob{ID: i, Priority: 2, Data: fmt.Sprintf("normal-%d", i)}
		}
		close(lowPriorityJobs)
	}()

	// Collect results
	time.Sleep(1 * time.Second)
	close(results)

	for result := range results {
		fmt.Printf("Result: %s\n", result)
	}
}

// WorkerPoolWithBatching demonstrates worker pool with batch processing
func WorkerPoolWithBatching() {
	fmt.Println("\n=== Worker Pool with Batching ===")

	const numWorkers = 2
	const batchSize = 3

	jobs := make(chan int, 10)
	results := make(chan []int, 10)

	// Start workers that process batches
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			batch := make([]int, 0, batchSize)

			for job := range jobs {
				batch = append(batch, job)

				if len(batch) == batchSize {
					fmt.Printf("Worker %d processing batch: %v\n", id, batch)
					time.Sleep(200 * time.Millisecond) // Simulate batch processing
					results <- batch
					batch = make([]int, 0, batchSize)
				}
			}

			// Process remaining items
			if len(batch) > 0 {
				fmt.Printf("Worker %d processing final batch: %v\n", id, batch)
				time.Sleep(200 * time.Millisecond)
				results <- batch
			}
		}(w)
	}

	// Send jobs
	for j := 1; j <= 7; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect batch results
	for r := 1; r <= 3; r++ {
		batch := <-results
		fmt.Printf("Batch result: %v\n", batch)
	}
}

// WorkerPoolWithMetrics demonstrates worker pool with metrics collection
func WorkerPoolWithMetrics() {
	fmt.Println("\n=== Worker Pool with Metrics ===")

	type Metrics struct {
		JobsProcessed int
		TotalTime     time.Duration
		Errors        int
	}

	const numWorkers = 3
	const numJobs = 9

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	metrics := make(chan Metrics, numWorkers)

	// Start workers with metrics
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			start := time.Now()
			processed := 0
			errors := 0

			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", id, job)

				// Simulate work
				time.Sleep(100 * time.Millisecond)

				// Simulate occasional errors
				if job%4 == 0 {
					errors++
					fmt.Printf("Worker %d: job %d failed\n", id, job)
				} else {
					results <- job * 5
					processed++
				}
			}

			// Send metrics
			metrics <- Metrics{
				JobsProcessed: processed,
				TotalTime:     time.Since(start),
				Errors:        errors,
			}
		}(w)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for r := 1; r <= numJobs; r++ {
		select {
		case result := <-results:
			fmt.Printf("Result: %d\n", result)
		case <-time.After(100 * time.Millisecond):
			// Handle timeout for failed jobs
		}
	}

	// Collect metrics
	close(results)
	for w := 1; w <= numWorkers; w++ {
		m := <-metrics
		fmt.Printf("Worker %d metrics: %+v\n", w, m)
	}
}

// WorkerPoolWithDynamicScaling demonstrates worker pool with dynamic scaling
func WorkerPoolWithDynamicScaling() {
	fmt.Println("\n=== Worker Pool with Dynamic Scaling ===")

	const maxWorkers = 5
	const numJobs = 15

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	workerCount := 0
	var mu sync.Mutex

	// Function to start a worker
	startWorker := func(id int) {
		mu.Lock()
		workerCount++
		mu.Unlock()

		defer func() {
			mu.Lock()
			workerCount--
			mu.Unlock()
		}()

		for job := range jobs {
			fmt.Printf("Worker %d processing job %d (total workers: %d)\n", id, job, workerCount)
			time.Sleep(100 * time.Millisecond)
			results <- job * 6
		}
	}

	// Start initial workers
	for w := 1; w <= 2; w++ {
		go startWorker(w)
	}

	// Send jobs and scale workers based on load
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j

			// Scale up if needed
			mu.Lock()
			if workerCount < maxWorkers && j%3 == 0 {
				go startWorker(workerCount + 1)
			}
			mu.Unlock()

			time.Sleep(50 * time.Millisecond)
		}
		close(jobs)
	}()

	// Collect results
	for r := 1; r <= numJobs; r++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}

// RunAllWorkerPoolExamples runs all worker pool examples
func RunAllWorkerPoolExamples() {
	fmt.Println("Running Worker Pool Examples...")

	//BasicWorkerPool()
	//WorkerPoolWithWaitGroup()
	//WorkerPoolWithContext()
	//WorkerPoolWithErrorHandling()
	//WorkerPoolWithRateLimiting()
	//WorkerPoolWithPriority()
	//WorkerPoolWithBatching()
	//WorkerPoolWithMetrics()
	WorkerPoolWithDynamicScaling()

	fmt.Println("\n=== All Worker Pool Examples Completed ===")
}
