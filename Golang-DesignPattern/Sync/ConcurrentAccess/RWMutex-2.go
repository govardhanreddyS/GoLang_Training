package main

import (
	"fmt"
	"sync"
)

// Counter represents a counter with safe concurrent access
type Counter struct {
	value int
	mu    sync.RWMutex
}

// NewCounter creates a new Counter instance
func NewCounter() *Counter {
	return &Counter{}
}

// Increment increments the counter value
func (c *Counter) Increment() {
	c.mu.Lock()         // Acquire a write lock
	defer c.mu.Unlock() // Ensure the lock is released when done
	c.value++
	fmt.Println("Counter value:", c.value) // Print the counter value immediately after incrementing
}

// Value returns the current value of the counter
func (c *Counter) Value() int {
	c.mu.RLock()         // Acquire a read lock
	defer c.mu.RUnlock() // Ensure the lock is released when done
	return c.value
}

func main() {
	counter := NewCounter()

	// Goroutines to increment the counter concurrently
	const numIncrementers = 10
	const numIncrements = 3
	var wg sync.WaitGroup
	wg.Add(numIncrementers)

	for i := 0; i < numIncrementers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				counter.Increment() // Increment the counter value
			}
		}()
	}

	wg.Wait() // Wait for all incrementer goroutines to finish
}
