package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var lock sync.Mutex
	cond := sync.NewCond(&lock)

	// WaitGroup to ensure all goroutines are done before main exits
	var wg sync.WaitGroup

	// Goroutine 1 waiting for the wake-up call
	wg.Add(1)
	go func() {
		defer wg.Done()

		lock.Lock()
		defer lock.Unlock()

		fmt.Println("Goroutine 1 is waiting for a wake-up call")
		cond.Wait()
		fmt.Println("Goroutine 1 received the wake-up call")
	}()

	// Goroutine 2 waiting for the wake-up call
	wg.Add(1)
	go func() {
		defer wg.Done()

		lock.Lock()
		defer lock.Unlock()

		fmt.Println("Goroutine 2 is waiting for a wake-up call")
		cond.Wait()
		fmt.Println("Goroutine 2 received the wake-up call")
	}()

	// Goroutine 3 waiting for the wake-up call
	wg.Add(1)
	go func() {
		defer wg.Done()

		lock.Lock()
		defer lock.Unlock()

		fmt.Println("Goroutine 3 is waiting for a wake-up call")
		cond.Wait()
		fmt.Println("Goroutine 3 received the wake-up call")
	}()

	// Goroutine responsible for broadcasting the message
	go func() {
		// Simulating some condition to be met
		fmt.Println("Waiting for condition to be met...")
		time.Sleep(2 * time.Second)

		// Send wake-up call to other goroutines
		lock.Lock()
		defer lock.Unlock()

		fmt.Println("Broadcasting message to wake up other goroutines")
		cond.Broadcast()
	}()

	// Wait for all goroutines to finish before exiting
	wg.Wait()

	fmt.Println("All goroutines have completed execution. Exiting main.")
}
