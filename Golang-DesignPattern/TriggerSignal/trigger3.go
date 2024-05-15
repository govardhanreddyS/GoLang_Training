package main

import (
	"fmt"
	"time"
)

func worker(id int, done chan bool) {
	fmt.Printf("Worker %d: Waiting for signal...\n", id)
	<-done // Wait for a signal on the channel
	fmt.Printf("Worker %d: Got signal. Continuing...\n", id)
}

func main() {
	numWorkers := 5
	doneChannels := make([]chan bool, numWorkers) // Create a pool of channels

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		doneChannels[i] = make(chan bool) // Create a channel for each worker
		go worker(i+1, doneChannels[i])   // Start the worker goroutines
	}

	// Simulate some work in the main goroutine
	time.Sleep(1 * time.Second)

	// Send signals to unblock specific workers based on conditions
	for i := 0; i < numWorkers; i++ {
		time.Sleep(500 * time.Millisecond)
		doneChannels[i] <- true // Signal to unblock the worker
	}

	// Give some time for the workers to finish
	time.Sleep(1 * time.Second)
}
