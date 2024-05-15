package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    value int
    mutex sync.RWMutex
)

func readWithRWMutex(id int) {
    mutex.RLock()
    defer mutex.RUnlock()
    fmt.Printf("Reader %d: Reading value with RWMutex: %d\n", id, value)
}

func writeWithRWMutex(id int, newValue int) {
    mutex.Lock()
    defer mutex.Unlock()
    fmt.Printf("Writer %d: Writing value with RWMutex\n", id)
    value = newValue
}

func readWithMutex(id int) {
    mutex.Lock()
    defer mutex.Unlock()
    fmt.Printf("Reader %d: Reading value with Mutex: %d\n", id, value)
}

func writeWithMutex(id int, newValue int) {
    mutex.Lock()
    defer mutex.Unlock()
    fmt.Printf("Writer %d: Writing value with Mutex\n", id)
    value = newValue
}

func main() {
    // Using RWMutex
    for i := 1; i <= 5; i++ {
        go readWithRWMutex(i)
    }
    time.Sleep(500 * time.Millisecond)
    go writeWithRWMutex(1, 10)
    time.Sleep(500 * time.Millisecond)
    for i := 6; i <= 10; i++ {
        go readWithRWMutex(i)
    }

    time.Sleep(2 * time.Second)

    fmt.Println("----")

    // Using Mutex
    for i := 1; i <= 5; i++ {
        go readWithMutex(i)
    }
    time.Sleep(500 * time.Millisecond)
    go writeWithMutex(1, 20)
    time.Sleep(500 * time.Millisecond)
    for i := 6; i <= 10; i++ {
        go readWithMutex(i)
    }

    time.Sleep(2 * time.Second)
}
