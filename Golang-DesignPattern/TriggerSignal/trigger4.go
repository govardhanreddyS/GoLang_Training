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

func otherRoutine(id int, done chan bool) {
	fmt.Printf("OtherRoutine %d: Waiting for signal...\n", id)
	<-done // Wait for a signal on the channel
	fmt.Printf("OtherRoutine %d: Got signal. Continuing...\n", id)
}

func main() {
	numWorkers := 3
	numOtherRoutines := 2

	doneChannels := make([]chan bool, numWorkers+numOtherRoutines) // Create a pool of channels

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		doneChannels[i] = make(chan bool) // Create a channel for each worker
		go worker(i+1, doneChannels[i])   // Start the worker goroutines
	}

	// Start other routines
	for i := 0; i < numOtherRoutines; i++ {
		doneChannels[numWorkers+i] = make(chan bool) // Create a channel for each other routine
		go otherRoutine(i+1, doneChannels[numWorkers+i]) // Start the other routines
	}

	// Simulate some work in the main goroutine
	time.Sleep(1 * time.Second)

	// Send signals to unblock specific workers or other routines based on conditions
	for i := 0; i < numWorkers+numOtherRoutines; i++ {
		time.Sleep(500 * time.Millisecond)
		doneChannels[i] <- true // Signal to unblock the worker or other routine
	}

	// Give some time for the workers and other routines to finish
	time.Sleep(1 * time.Second)
}
