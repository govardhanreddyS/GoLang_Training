package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        time.Sleep(2 * time.Second)
        cancel()
    }()

    select {
    case <-time.After(5 * time.Second):
        fmt.Println("Operation completed successfully")
    case <-ctx.Done():
        fmt.Println("Operation cancelled")
    }
}
