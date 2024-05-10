package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Worker finished task.")
	case <-ctx.Done():
		fmt.Println("Worker received cancellation signal.")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel() // set the value for Done

	go worker(ctx)

	// Simulate cancellation after 2 seconds
	time.Sleep(2 * time.Second)

	fmt.Println("Main function exiting.")
}

