package main

import (
    "fmt"
    "sync"
)

var (
    counter = 0
    mutex   sync.Mutex
)

func increment() {
    mutex.Lock()
    defer mutex.Unlock()
    counter++
    fmt.Println("Incremented counter to:", counter)
}

func main() {
    for i := 0; i < 10; i++ {
        go increment()
    }
    // Wait for goroutines to finish
    fmt.Scanln()
}
