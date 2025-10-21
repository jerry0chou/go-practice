package concurrency

import (
	"fmt"
	"sync"
	"time"
)

// BasicMutex demonstrates basic mutex usage for protecting shared state
func BasicMutex() {
	fmt.Println("=== Basic Mutex ===")

	var mu sync.Mutex
	counter := 0

	// Start multiple goroutines that increment counter
	for i := 1; i <= 5; i++ {
		go func(id int) {
			mu.Lock()
			defer mu.Unlock()

			// Critical section
			oldValue := counter
			time.Sleep(10 * time.Millisecond) // Simulate some work
			counter = oldValue + 1
			fmt.Printf("Goroutine %d: counter = %d\n", id, counter)
		}(i)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Final counter value: %d\n", counter)
}

// MutexWithoutProtection demonstrates what happens without mutex protection
func MutexWithoutProtection() {
	fmt.Println("\n=== Without Mutex Protection (Race Condition) ===")

	counter := 0

	// Start multiple goroutines without protection
	for i := 1; i <= 5; i++ {
		go func(id int) {
			// Race condition: multiple goroutines accessing counter simultaneously
			oldValue := counter
			time.Sleep(10 * time.Millisecond)
			counter = oldValue + 1
			fmt.Printf("Goroutine %d: counter = %d\n", id, counter)
		}(i)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Final counter value (may be incorrect): %d\n", counter)
}

// RWMutexBasic demonstrates basic RWMutex usage
func RWMutexBasic() {
	fmt.Println("\n=== Basic RWMutex ===")

	var rwmu sync.RWMutex
	data := make(map[string]int)

	// Start writer goroutines
	for i := 1; i <= 3; i++ {
		go func(id int) {
			rwmu.Lock()
			defer rwmu.Unlock()

			key := fmt.Sprintf("key%d", id)
			data[key] = id * 10
			fmt.Printf("Writer %d: wrote %s = %d\n", id, key, data[key])
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	// Start reader goroutines
	for i := 1; i <= 5; i++ {
		go func(id int) {
			rwmu.RLock()
			defer rwmu.RUnlock()

			// Multiple readers can access simultaneously
			fmt.Printf("Reader %d: data has %d entries\n", id, len(data))
			time.Sleep(30 * time.Millisecond)
		}(i)
	}

	time.Sleep(300 * time.Millisecond)
	fmt.Printf("Final data: %v\n", data)
}

// MutexWithDefer demonstrates proper defer usage with mutexes
func MutexWithDefer() {
	fmt.Println("\n=== Mutex with Defer ===")

	var mu sync.Mutex
	sharedResource := "initial"

	// Start goroutines with proper defer usage
	for i := 1; i <= 3; i++ {
		go func(id int) {
			mu.Lock()
			defer mu.Unlock() // Ensures unlock even if panic occurs

			fmt.Printf("Goroutine %d: accessing shared resource\n", id)
			sharedResource = fmt.Sprintf("modified by %d", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Goroutine %d: resource = %s\n", id, sharedResource)
		}(i)
	}

	time.Sleep(400 * time.Millisecond)
	fmt.Printf("Final shared resource: %s\n", sharedResource)
}

// MutexWithTryLock demonstrates try-lock functionality (Go 1.18+)
func MutexWithTryLock() {
	fmt.Println("\n=== Mutex with TryLock ===")

	var mu sync.Mutex
	counter := 0

	// Start goroutines that try to acquire lock
	for i := 1; i <= 3; i++ {
		go func(id int) {
			// Try to acquire lock without blocking
			if mu.TryLock() {
				defer mu.Unlock()

				fmt.Printf("Goroutine %d: acquired lock\n", id)
				counter++
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("Goroutine %d: counter = %d\n", id, counter)
			} else {
				fmt.Printf("Goroutine %d: could not acquire lock\n", id)
			}
		}(i)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Final counter value: %d\n", counter)
}

// RWMutexWithMultipleReaders demonstrates multiple readers with RWMutex
func RWMutexWithMultipleReaders() {
	fmt.Println("\n=== RWMutex with Multiple Readers ===")

	var rwmu sync.RWMutex
	sharedData := make(map[string]string)

	// Start multiple readers
	for i := 1; i <= 5; i++ {
		go func(id int) {
			rwmu.RLock()
			defer rwmu.RUnlock()

			fmt.Printf("Reader %d: reading data\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Reader %d: data has %d entries\n", id, len(sharedData))
		}(i)
	}

	// Start a writer
	go func() {
		rwmu.Lock()
		defer rwmu.Unlock()

		fmt.Println("Writer: updating data")
		sharedData["key1"] = "value1"
		sharedData["key2"] = "value2"
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Writer: data updated")
	}()

	time.Sleep(300 * time.Millisecond)
	fmt.Printf("Final data: %v\n", sharedData)
}

// MutexWithTimeout demonstrates timeout with mutex (using channels)
func MutexWithTimeout() {
	fmt.Println("\n=== Mutex with Timeout ===")

	var mu sync.Mutex
	acquired := make(chan bool, 1)

	// First goroutine holds the lock
	go func() {
		mu.Lock()
		fmt.Println("First goroutine: acquired lock")
		time.Sleep(200 * time.Millisecond)
		mu.Unlock()
		fmt.Println("First goroutine: released lock")
	}()

	// Second goroutine tries to acquire with timeout
	go func() {
		time.Sleep(50 * time.Millisecond) // Wait a bit

		select {
		case <-time.After(100 * time.Millisecond):
			fmt.Println("Second goroutine: timeout waiting for lock")
		default:
			mu.Lock()
			defer mu.Unlock()
			fmt.Println("Second goroutine: acquired lock")
			acquired <- true
		}
	}()

	time.Sleep(300 * time.Millisecond)

	select {
	case <-acquired:
		fmt.Println("Lock was acquired successfully")
	default:
		fmt.Println("Lock acquisition timed out")
	}
}

// MutexWithConditionalAccess demonstrates conditional access patterns
func MutexWithConditionalAccess() {
	fmt.Println("\n=== Mutex with Conditional Access ===")

	var mu sync.Mutex
	condition := false
	waiting := 0

	// Start goroutines that wait for condition
	for i := 1; i <= 3; i++ {
		go func(id int) {
			mu.Lock()
			waiting++
			fmt.Printf("Goroutine %d: waiting for condition (waiting: %d)\n", id, waiting)

			// Wait for condition (simplified - in real code use sync.Cond)
			for !condition {
				mu.Unlock()
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
			}

			waiting--
			fmt.Printf("Goroutine %d: condition met, proceeding\n", id)
			mu.Unlock()
		}(i)
	}

	// Set condition after delay
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	condition = true
	fmt.Println("Condition set to true")
	mu.Unlock()

	time.Sleep(200 * time.Millisecond)
}

// MutexWithDeadlockPrevention demonstrates deadlock prevention
func MutexWithDeadlockPrevention() {
	fmt.Println("\n=== Mutex with Deadlock Prevention ===")

	var mu1, mu2 sync.Mutex

	// Function to acquire locks in consistent order
	acquireLocks := func(id int, mu1, mu2 *sync.Mutex) {
		// Always acquire mu1 first, then mu2
		mu1.Lock()
		defer mu1.Unlock()

		fmt.Printf("Goroutine %d: acquired mu1\n", id)
		time.Sleep(50 * time.Millisecond)

		mu2.Lock()
		defer mu2.Unlock()

		fmt.Printf("Goroutine %d: acquired mu2\n", id)
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("Goroutine %d: completed work\n", id)
	}

	// Start goroutines with consistent lock ordering
	for i := 1; i <= 3; i++ {
		go acquireLocks(i, &mu1, &mu2)
	}

	time.Sleep(300 * time.Millisecond)
	fmt.Println("All goroutines completed without deadlock")
}

// MutexWithResourcePool demonstrates mutex with resource pool
func MutexWithResourcePool() {
	fmt.Println("\n=== Mutex with Resource Pool ===")

	var mu sync.Mutex
	available := make([]int, 0)
	inUse := make([]int, 0)

	// Initialize resource pool
	for i := 1; i <= 3; i++ {
		available = append(available, i)
	}

	// Function to acquire resource
	acquireResource := func(id int) (int, bool) {
		mu.Lock()
		defer mu.Unlock()

		if len(available) > 0 {
			resource := available[0]
			available = available[1:]
			inUse = append(inUse, resource)
			fmt.Printf("Goroutine %d: acquired resource %d\n", id, resource)
			return resource, true
		}
		fmt.Printf("Goroutine %d: no resources available\n", id)
		return 0, false
	}

	// Function to release resource
	releaseResource := func(id, resource int) {
		mu.Lock()
		defer mu.Unlock()

		// Remove from inUse
		for i, r := range inUse {
			if r == resource {
				inUse = append(inUse[:i], inUse[i+1:]...)
				break
			}
		}

		// Add back to available
		available = append(available, resource)
		fmt.Printf("Goroutine %d: released resource %d\n", id, resource)
	}

	// Start goroutines that use resources
	for i := 1; i <= 5; i++ {
		go func(id int) {
			resource, ok := acquireResource(id)
			if ok {
				time.Sleep(100 * time.Millisecond)
				releaseResource(id, resource)
			}
		}(i)
	}

	time.Sleep(300 * time.Millisecond)
	fmt.Printf("Final state - Available: %v, InUse: %v\n", available, inUse)
}

// RunAllMutexExamples runs all mutex examples
func RunAllMutexExamples() {
	fmt.Println("Running Mutex Examples...")

	BasicMutex()
	MutexWithoutProtection()
	RWMutexBasic()
	MutexWithDefer()
	MutexWithTryLock()
	RWMutexWithMultipleReaders()
	MutexWithTimeout()
	MutexWithConditionalAccess()
	MutexWithDeadlockPrevention()
	MutexWithResourcePool()

	fmt.Println("\n=== All Mutex Examples Completed ===")
}
