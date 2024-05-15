package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    deadline := time.Now().Add(2 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()

    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Operation completed successfully")
    case <-ctx.Done():
        fmt.Println("Operation deadline exceeded")
    }
}
