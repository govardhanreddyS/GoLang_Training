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

func readWithRWMutex() {
    mutex.RLock()
    defer mutex.RUnlock()
    fmt.Println("Reading value with RWMutex:", value)
}

func writeWithRWMutex(newValue int) {
    mutex.Lock()
    defer mutex.Unlock()
    fmt.Println("Writing value with RWMutex")
    value = newValue
}

func readWithMutex() {
    mutex.Lock()
    defer mutex.Unlock()
    fmt.Println("Reading value with Mutex:", value)
}

func writeWithMutex(newValue int) {
    mutex.Lock()
    defer mutex.Unlock()
    fmt.Println("Writing value with Mutex")
    value = newValue
}

func main() {
    // Using RWMutex
    go readWithRWMutex()
    go writeWithRWMutex(10)
    go readWithRWMutex()

    time.Sleep(1 * time.Second)
    
    fmt.Println("----")

    // Using Mutex
    go readWithMutex()
    go writeWithMutex(20)
    go readWithMutex()

    time.Sleep(1 * time.Second)
}
