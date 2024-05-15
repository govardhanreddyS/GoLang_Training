package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context) {
    select {
    case <-ctx.Done():
        fmt.Println("Worker received cancellation signal")
        return
    case <-time.After(2* time.Second):
        fmt.Println("Worker completed task")
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    go worker(ctx)

    // Simulate cancellation after 1 second
    time.Sleep(1 * time.Second)
    cancel()
    time.Sleep(1 * time.Second) // Allow time for worker to finish
}
