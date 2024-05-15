package main

import (
	"fmt"
	"sync"
)

func main() {
	// Create channels for synchronization
	channels := make([]chan bool, 5)
	for i := range channels {
		channels[i] = make(chan bool)
	}

	// WaitGroup to ensure all goroutines finish before main exits
	var wg sync.WaitGroup
	wg.Add(5)

	// Goroutines to print messages in a synchronized manner
	for i := range channels {
		go func(id int, ch chan bool) {
			defer wg.Done()
			for {
				select {
				case <-ch:
					fmt.Printf("Goroutine %d: Executing\n", id)
					// Simulate some work
					// Here, we sleep for a short duration to simulate work
					// You can replace this with your actual task
					// For demonstration purpose, let's sleep for 500 milliseconds
					// You may adjust the duration according to your needs
					//time.Sleep(500 * time.Millisecond)
					// After execution, signal the next goroutine to execute
					channels[(id+1)%5] <- true
				}
			}
		}(i, channels[i])
	}

	// Start the execution by signaling the first goroutine
	channels[0] <- true

	// Wait for all goroutines to finish
	wg.Wait()
}
