package main

import (
	"fmt"
	"sync"
)

type Resource struct {
	mu sync.Mutex
}

func main() {
	resA := &Resource{}
	resB := &Resource{}

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1
	go func() {
		defer wg.Done()

		resA.mu.Lock()
		defer resA.mu.Unlock()

		// Simulate some processing time
		fmt.Println("Goroutine 1 acquired resource A")
		fmt.Println("Goroutine 1 waiting to acquire resource B")
		resB.mu.Lock()
		defer resB.mu.Unlock()

		fmt.Println("Goroutine 1 acquired resource B")
	}()

	// Goroutine 2
	go func() {
		defer wg.Done()

		resB.mu.Lock()
		defer resB.mu.Unlock()

		// Simulate some processing time
		fmt.Println("Goroutine 2 acquired resource B")
		fmt.Println("Goroutine 2 waiting to acquire resource A")
		resA.mu.Lock()
		defer resA.mu.Unlock()

		fmt.Println("Goroutine 2 acquired resource A")
	}()

	wg.Wait()
	fmt.Println("Both goroutines finished successfully")
}
