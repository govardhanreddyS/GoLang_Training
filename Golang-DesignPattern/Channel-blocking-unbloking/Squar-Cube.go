package main

import (
	"fmt"
	"sync"
)

func main() {
	// Create channels for communication between producer and consumers
	squareCh := make(chan int)
	cubeCh := make(chan int)
	squareDone := make(chan struct{})
	cubeDone := make(chan struct{})

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Add 2 to the WaitGroup for the squaring and cubing goroutines
	wg.Add(2)

	// Producer goroutine
	go func() {
		defer close(squareCh)
		defer close(cubeCh)

		for i := 1; i <= 10; i++ {
			squareCh <- i
			cubeCh <- i
		}
	}()

	// Squaring goroutine
	go func() {
		defer wg.Done()

		for val := range squareCh {
			fmt.Println("Receiver 1 squaring the value:", val, "->", val*val)
			cubeDone <- struct{}{} // Signal that squaring is done for this value
			<-squareDone          // Wait for a signal from the cubing goroutine
		}
	}()

	// Cubing goroutine
	go func() {
		defer wg.Done()

		for val := range cubeCh {
			<-cubeDone                     // Wait for a signal from the squaring goroutine
			fmt.Println("Receiver 2 cubing the value:", val, "->", val*val*val)
			squareDone <- struct{}{} // Signal that cubing is done for this value
		}
	}()

	// Wait for all squaring and cubing operations to complete
	wg.Wait()
}
