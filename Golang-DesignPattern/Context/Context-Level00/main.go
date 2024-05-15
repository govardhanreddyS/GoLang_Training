package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker received cancellation signal.")
			return
		default:
			fmt.Println("Worker is still working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go worker(ctx)

	// Simulate cancellation after 3 seconds
	time.Sleep(3 * time.Second)
	cancel()

	// Wait for the worker to gracefully exit
	time.Sleep(1 * time.Second)
	fmt.Println("Main function exiting.")
}
