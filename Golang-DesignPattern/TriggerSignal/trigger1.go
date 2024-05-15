package main

import (
	"fmt"
	"time"
)

func worker(done chan bool) {
	fmt.Println("Worker: Doing some work...")
	time.Sleep(2 * time.Second) // Simulate some work
	fmt.Println("Worker: Work done.")
	done <- true // Signal that work is done by sending a value to the channel
}

func main() {
	done := make(chan bool) // Create a channel

	go worker(done) // Start the worker goroutine

	// Block until receiving a signal from the worker goroutine
	<-done
	fmt.Println("Main: Worker has finished its work. Proceeding...")
}
