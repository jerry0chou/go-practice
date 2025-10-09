package data_structure

import (
	"container/heap"
	"container/list"
	"container/ring"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ListOperations demonstrates various operations with container/list
func ListOperations() {
	fmt.Println("=== List Operations (container/list) ===")

	l := list.New()
	fmt.Println("✅ Created new list")
	l.PushBack("first")
	l.PushBack("second")
	l.PushFront("zero")
	l.PushBack("third")
	fmt.Println("✅ Added elements to list")

	fmt.Println("📋 List contents (front to back):")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("  %v\n", e.Value)
	}

	fmt.Println("📋 List contents (back to front):")
	for e := l.Back(); e != nil; e = e.Prev() {
		fmt.Printf("  %v\n", e.Value)
	}

	secondElement := l.Front().Next().Next() // "second"
	l.InsertAfter("inserted_after_second", secondElement)
	fmt.Println("✅ Inserted element after 'second'")

	l.InsertBefore("inserted_before_second", secondElement)
	fmt.Println("✅ Inserted element before 'second'")

	l.Remove(secondElement)
	fmt.Println("✅ Removed 'second' element")

	fmt.Println("📋 Updated list contents:")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("  %v\n", e.Value)
	}

	lastElement := l.Back()
	l.MoveToFront(lastElement)
	fmt.Println("✅ Moved last element to front")

	firstElement := l.Front()
	l.MoveToBack(firstElement)
	fmt.Println("✅ Moved first element to back")

	fmt.Println("📋 Final list contents:")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("  %v\n", e.Value)
	}

	fmt.Printf("📊 List length: %d\n", l.Len())
}

// ListAdvancedOperations demonstrates advanced list operations
func ListAdvancedOperations() {
	fmt.Println("\n=== Advanced List Operations ===")

	// Create list with numbers
	l := list.New()
	for i := 1; i <= 5; i++ {
		l.PushBack(i)
	}

	fmt.Println("📋 Original list:")
	printList(l)

	// Find and remove specific value
	valueToRemove := 3
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == valueToRemove {
			l.Remove(e)
			fmt.Printf("✅ Removed value: %d\n", valueToRemove)
			break
		}
	}

	// Insert in sorted position
	newValue := 4
	inserted := false
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(int) > newValue {
			l.InsertBefore(newValue, e)
			inserted = true
			break
		}
	}
	if !inserted {
		l.PushBack(newValue)
	}
	fmt.Printf("✅ Inserted value in sorted position: %d\n", newValue)

	fmt.Println("📋 Updated list:")
	printList(l)
}

func printList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("  %v\n", e.Value)
	}
}

// GenericHeap - using generics with built-in heap
// This eliminates code duplication and provides type safety
type GenericHeap[T comparable] []T

func (h GenericHeap[T]) Len() int      { return len(h) }
func (h GenericHeap[T]) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *GenericHeap[T]) Push(x interface{}) {
	*h = append(*h, x.(T))
}

func (h *GenericHeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// IntHeap - min-heap for integers using generics
type IntHeap struct {
	GenericHeap[int]
}

func (h IntHeap) Less(i, j int) bool { return h.GenericHeap[i] < h.GenericHeap[j] }

// StringHeap - min-heap for strings using generics
type StringHeap struct {
	GenericHeap[string]
}

func (h StringHeap) Less(i, j int) bool { return h.GenericHeap[i] < h.GenericHeap[j] }

// Task represents a task with priority
type Task struct {
	Name     string
	Priority int
}

// TaskHeap - max-heap for tasks using generics
type TaskHeap struct {
	GenericHeap[Task]
}

func (h TaskHeap) Less(i, j int) bool { return h.GenericHeap[i].Priority > h.GenericHeap[j].Priority }

// HeapOperations demonstrates using built-in container/heap package
func HeapOperations() {
	fmt.Println("\n=== Heap Operations (container/heap) ===")

	fmt.Println("📊 Integer Min-Heap using generics with container/heap:")
	intHeap := &IntHeap{GenericHeap[int]{2, 1, 5, 3, 4}}
	heap.Init(intHeap)

	fmt.Printf("Initial heap: %v\n", intHeap.GenericHeap)
	fmt.Printf("Min element (root): %d\n", intHeap.GenericHeap[0])

	heap.Push(intHeap, 0)
	fmt.Printf("After pushing 0: %v\n", intHeap.GenericHeap)

	fmt.Println("Popping elements (min-heap order):")
	for intHeap.Len() > 0 {
		min := heap.Pop(intHeap).(int)
		fmt.Printf("  Popped: %d, Remaining: %v\n", min, intHeap.GenericHeap)
	}

	fmt.Println("\n📊 String Heap using generics with container/heap:")
	stringHeap := &StringHeap{GenericHeap[string]{"banana", "apple", "cherry", "date"}}
	heap.Init(stringHeap)

	fmt.Printf("Initial string heap: %v\n", stringHeap.GenericHeap)
	fmt.Printf("Min string (root): %s\n", stringHeap.GenericHeap[0])

	heap.Push(stringHeap, "apricot")
	fmt.Printf("After pushing 'apricot': %v\n", stringHeap.GenericHeap)

	fmt.Println("Popping strings (alphabetical order):")
	for stringHeap.Len() > 0 {
		str := heap.Pop(stringHeap).(string)
		fmt.Printf("  Popped: %s\n", str)
	}

	fmt.Println("\n📊 Task Priority Queue using generics with container/heap:")
	taskHeap := &TaskHeap{
		GenericHeap[Task]{
			{"Write docs", 2},
			{"Fix bug", 5},
			{"Code review", 3},
			{"Deploy", 1},
		},
	}
	heap.Init(taskHeap)

	fmt.Println("Initial task heap:")
	for _, task := range taskHeap.GenericHeap {
		fmt.Printf("  %s (priority: %d)\n", task.Name, task.Priority)
	}

	heap.Push(taskHeap, Task{"Emergency fix", 10})
	heap.Push(taskHeap, Task{"Update README", 1})

	fmt.Println("\nProcessing tasks by priority (highest first):")
	for taskHeap.Len() > 0 {
		task := heap.Pop(taskHeap).(Task)
		fmt.Printf("  Processing: %s (priority: %d)\n", task.Name, task.Priority)
	}

	fmt.Println("\n📊 Heap.Fix for updating priorities:")

	updateHeap := &TaskHeap{
		GenericHeap[Task]{
			{"Task A", 3},
			{"Task B", 1},
			{"Task C", 5},
		},
	}
	heap.Init(updateHeap)

	fmt.Println("Initial heap:")
	for i, task := range updateHeap.GenericHeap {
		fmt.Printf("  [%d] %s (priority: %d)\n", i, task.Name, task.Priority)
	}

	updateHeap.GenericHeap[1].Priority = 8
	heap.Fix(updateHeap, 1)

	fmt.Println("After updating Task B priority to 8:")
	for i, task := range updateHeap.GenericHeap {
		fmt.Printf("  [%d] %s (priority: %d)\n", i, task.Name, task.Priority)
	}

	fmt.Println("Popping updated heap:")
	for updateHeap.Len() > 0 {
		task := heap.Pop(updateHeap).(Task)
		fmt.Printf("  %s (priority: %d)\n", task.Name, task.Priority)
	}
}

// RingOperations demonstrates ring operations
func RingOperations() {
	fmt.Println("\n=== Ring Operations (container/ring) ===")

	// Create a ring with 5 elements
	r := ring.New(5)
	fmt.Println("✅ Created ring with 5 elements")

	// Initialize ring with values
	for i := 0; i < 5; i++ {
		r.Value = i
		r = r.Next()
	}

	// Display ring contents
	fmt.Println("📋 Ring contents:")
	printRing(r, 5)

	// Move to next element
	r = r.Next()
	fmt.Printf("✅ Moved to next element: %d\n", r.Value)

	// Move to previous element
	r = r.Prev()
	fmt.Printf("✅ Moved to previous element: %d\n", r.Value)

	// Move multiple steps
	r = r.Move(3)
	fmt.Printf("✅ Moved 3 steps forward: %d\n", r.Value)

	// Link rings
	fmt.Println("\n📊 Ring Linking Example:")
	r1 := ring.New(3)
	r2 := ring.New(2)

	// Initialize first ring
	for i := 0; i < 3; i++ {
		r1.Value = fmt.Sprintf("R1-%d", i)
		r1 = r1.Next()
	}

	// Initialize second ring
	for i := 0; i < 2; i++ {
		r2.Value = fmt.Sprintf("R2-%d", i)
		r2 = r2.Next()
	}

	fmt.Println("Ring 1:")
	printRing(r1, 3)
	fmt.Println("Ring 2:")
	printRing(r2, 2)

	// Link the rings
	r1.Link(r2)
	fmt.Println("After linking:")
	printRing(r1, 5)

	// Unlink elements
	unlinked := r1.Unlink(2)
	fmt.Println("After unlinking 2 elements:")
	printRing(r1, 3)
	fmt.Println("Unlinked elements:")
	printRing(unlinked, 2)
}

func printRing(r *ring.Ring, n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("  %v\n", r.Value)
		r = r.Next()
	}
}

// Person represents a person with name and age
type Person struct {
	Name string
	Age  int
}

// ByAge implements sort.Interface for []Person based on age
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

// ByName implements sort.Interface for []Person based on name
type ByName []Person

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SortOperations demonstrates various sorting operations
func SortOperations() {
	fmt.Println("\n=== Sort Operations (sort package) ===")

	// Basic integer sorting
	fmt.Println("📊 Integer Sorting:")
	numbers := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Original: %v\n", numbers)

	sort.Ints(numbers)
	fmt.Printf("Sorted: %v\n", numbers)

	// Check if sorted
	fmt.Printf("Is sorted: %t\n", sort.IntsAreSorted(numbers))

	// String sorting
	fmt.Println("\n📊 String Sorting:")
	names := []string{"Charlie", "Alice", "Bob", "David"}
	fmt.Printf("Original: %v\n", names)

	sort.Strings(names)
	fmt.Printf("Sorted: %v\n", names)

	// Float sorting
	fmt.Println("\n📊 Float Sorting:")
	floats := []float64{3.14, 2.71, 1.41, 1.73}
	fmt.Printf("Original: %v\n", floats)

	sort.Float64s(floats)
	fmt.Printf("Sorted: %v\n", floats)

	// Custom struct sorting
	fmt.Println("\n📊 Custom Struct Sorting:")
	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"David", 35},
	}

	fmt.Println("Original people:")
	printPeople(people)

	// Sort by age
	sort.Sort(ByAge(people))
	fmt.Println("Sorted by age:")
	printPeople(people)

	// Sort by name
	sort.Sort(ByName(people))
	fmt.Println("Sorted by name:")
	printPeople(people)

	// Reverse sorting
	fmt.Println("\n📊 Reverse Sorting:")
	numbers = []int{1, 2, 3, 4, 5}
	fmt.Printf("Original: %v\n", numbers)

	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	fmt.Printf("Reversed: %v\n", numbers)

	// Search operations
	fmt.Println("\n📊 Search Operations:")
	sortedNumbers := []int{1, 3, 5, 7, 9, 11, 13, 15}
	fmt.Printf("Sorted array: %v\n", sortedNumbers)

	// Binary search
	index := sort.SearchInts(sortedNumbers, 7)
	fmt.Printf("Index of 7: %d\n", index)

	index = sort.SearchInts(sortedNumbers, 6)
	fmt.Printf("Index where 6 would be inserted: %d\n", index)

	// Custom search
	index = sort.Search(len(sortedNumbers), func(i int) bool {
		return sortedNumbers[i] >= 10
	})
	fmt.Printf("Index of first element >= 10: %d\n", index)
}

func printPeople(people []Person) {
	for _, p := range people {
		fmt.Printf("  %s (age: %d)\n", p.Name, p.Age)
	}
}

// BuiltInPackageExamples demonstrates more built-in package features
func BuiltInPackageExamples() {
	fmt.Println("\n=== More Built-in Package Examples ===")

	// More sort package features
	fmt.Println("📊 Advanced Sort Package Features:")

	// Sort with custom comparison
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	fmt.Printf("Original words: %v\n", words)

	// Sort by length
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) < len(words[j])
	})
	fmt.Printf("Sorted by length: %v\n", words)

	// Sort by last character
	sort.Slice(words, func(i, j int) bool {
		return words[i][len(words[i])-1] < words[j][len(words[j])-1]
	})
	fmt.Printf("Sorted by last character: %v\n", words)

	// Stable sort example
	fmt.Println("\n📊 Stable Sort Example:")
	type Student struct {
		Name  string
		Grade int
	}

	students := []Student{
		{"Alice", 85},
		{"Bob", 90},
		{"Charlie", 85},
		{"David", 90},
		{"Eve", 85},
	}

	fmt.Println("Original students:")
	for _, s := range students {
		fmt.Printf("  %s: %d\n", s.Name, s.Grade)
	}

	// Stable sort by grade (preserves original order for equal grades)
	sort.SliceStable(students, func(i, j int) bool {
		return students[i].Grade > students[j].Grade
	})

	fmt.Println("Stable sorted by grade (descending):")
	for _, s := range students {
		fmt.Printf("  %s: %d\n", s.Name, s.Grade)
	}

	// More heap operations using built-in container/heap
	fmt.Println("\n📊 More Built-in Heap Operations:")

	// Max-heap using generics with built-in heap package
	realMaxHeap := &MaxIntHeap{}
	heap.Init(realMaxHeap)

	values := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("Creating max-heap using generics:")
	for _, v := range values {
		heap.Push(realMaxHeap, v)
		fmt.Printf("  Pushed %d, max-heap: %v\n", v, realMaxHeap.GenericHeap)
	}

	fmt.Println("Popping from max-heap (largest first):")
	for realMaxHeap.Len() > 0 {
		fmt.Printf("  Popped: %d\n", heap.Pop(realMaxHeap))
	}

	// More list operations with different data types
	fmt.Println("\n📊 List with Different Data Types:")

	// List of structs
	type Task struct {
		ID   int
		Name string
		Done bool
	}

	taskList := list.New()
	tasks := []Task{
		{1, "Write code", false},
		{2, "Test code", false},
		{3, "Document code", false},
		{4, "Deploy code", false},
	}

	for _, task := range tasks {
		taskList.PushBack(task)
	}

	fmt.Println("Task list:")
	for e := taskList.Front(); e != nil; e = e.Next() {
		task := e.Value.(Task)
		status := "❌"
		if task.Done {
			status = "✅"
		}
		fmt.Printf("  %s %d: %s\n", status, task.ID, task.Name)
	}

	// More ring operations - practical use case
	fmt.Println("\n📊 Ring for Circular Buffer (Sliding Window):")

	// Create a ring to store last 5 numbers
	window := ring.New(5)

	// Simulate incoming data stream
	dataStream := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Println("Processing data stream with sliding window of 5:")
	for i, data := range dataStream {
		window.Value = data
		window = window.Next()

		if i >= 4 { // Start showing window after we have 5 elements
			fmt.Printf("  Data %d: Window = [", data)
			window.Do(func(x interface{}) {
				if x != nil {
					fmt.Printf("%d ", x.(int))
				}
			})
			fmt.Println("]")
		}
	}

	// Ring for round-robin scheduling
	fmt.Println("\n📊 Ring for Round-Robin Scheduling:")

	processes := []string{"Process A", "Process B", "Process C"}
	processRing := ring.New(len(processes))

	// Initialize ring with processes
	for _, process := range processes {
		processRing.Value = process
		processRing = processRing.Next()
	}

	fmt.Println("Round-robin scheduling (3 cycles):")
	for cycle := 0; cycle < 3; cycle++ {
		fmt.Printf("  Cycle %d:\n", cycle+1)
		for i := 0; i < len(processes); i++ {
			fmt.Printf("    Executing: %s\n", processRing.Value)
			processRing = processRing.Next()
		}
	}
}

// MaxIntHeap is a max-heap of integers using generics
type MaxIntHeap struct {
	GenericHeap[int]
}

func (h MaxIntHeap) Less(i, j int) bool { return h.GenericHeap[i] > h.GenericHeap[j] }

// UtilityFunctions demonstrates various utility functions
func UtilityFunctions() {
	fmt.Println("\n=== Utility Functions ===")

	// String manipulation
	fmt.Println("📊 String Operations:")
	str := "Hello, World!"
	fmt.Printf("Original: %s\n", str)
	fmt.Printf("Length: %d\n", len(str))
	fmt.Printf("Uppercase: %s\n", strings.ToUpper(str))
	fmt.Printf("Lowercase: %s\n", strings.ToLower(str))
	fmt.Printf("Contains 'World': %t\n", strings.Contains(str, "World"))
	fmt.Printf("Split by comma: %v\n", strings.Split(str, ","))

	// Number conversion
	fmt.Println("\n📊 Number Conversion:")
	numStr := "123"
	num, err := strconv.Atoi(numStr)
	if err == nil {
		fmt.Printf("String '%s' to int: %d\n", numStr, num)
	}

	num = 456
	numStr = strconv.Itoa(num)
	fmt.Printf("Int %d to string: '%s'\n", num, numStr)

	// Array/slice operations
	fmt.Println("\n📊 Slice Operations:")
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("Original slice: %v\n", slice)
	fmt.Printf("Length: %d, Capacity: %d\n", len(slice), cap(slice))

	// Append elements
	slice = append(slice, 6, 7, 8)
	fmt.Printf("After append: %v\n", slice)

	// Slice operations
	subSlice := slice[2:5]
	fmt.Printf("Sub-slice [2:5]: %v\n", subSlice)

	// Copy slice
	copied := make([]int, len(slice))
	copy(copied, slice)
	fmt.Printf("Copied slice: %v\n", copied)
}

// RunAllDataStructureExamples runs all data structure examples
func RunAllDataStructureExamples() {
	fmt.Println("🎯 Go Built-in Data Structure Package Examples")
	fmt.Println("===============================================")

	ListOperations()
	ListAdvancedOperations()
	HeapOperations()
	RingOperations()
	SortOperations()
	BuiltInPackageExamples()
	UtilityFunctions()

	fmt.Println("\n✅ All built-in package operations completed!")
}
