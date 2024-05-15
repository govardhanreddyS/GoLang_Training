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

	// Goroutine responsible for broadcasting the message and waiting for others to complete
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Simulating some condition to be met
		fmt.Println("Waiting for condition to be met...")
		time.Sleep(2 * time.Second)

		// Send wake-up call to other goroutines
		lock.Lock()
		defer lock.Unlock()

		fmt.Println("Broadcasting message to wake up other goroutines")
		cond.Broadcast()
	}()

	// Goroutines waiting for the wake-up call
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			lock.Lock()
			defer lock.Unlock()

			fmt.Printf("Goroutine %d is waiting for a wake-up call\n", id)
			cond.Wait()
			fmt.Printf("Goroutine %d received the wake-up call\n", id)
		}(i)
	}

	// Wait for all goroutines to finish before exiting
	wg.Wait()

	fmt.Println("All goroutines have completed execution. Exiting main.")
}
