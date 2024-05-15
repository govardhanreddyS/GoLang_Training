package main

import (
    "fmt"
    "sync"
    "time"
)

type Data struct {
    Value int
}

type DataStore struct {
    data    Data
    rwMutex sync.RWMutex
}

func (ds *DataStore) Read() {
    ds.rwMutex.RLock()
    defer ds.rwMutex.RUnlock()

    fmt.Println("Reading data:", ds.data.Value)
    time.Sleep(time.Second)
}

func (ds *DataStore) Write(value int) {
    ds.rwMutex.Lock()
    defer ds.rwMutex.Unlock()

    fmt.Println("Writing data:", value)
    ds.data.Value = value
    time.Sleep(time.Second)
}

func main() {
    ds := &DataStore{}

    // Perform concurrent read operations
    for i := 0; i < 3; i++ {
        go func() {
            ds.Read()
        }()
    }

    // Perform concurrent write operations
    for i := 0; i < 2; i++ {
        go func() {
            newValue := i + 1
            ds.Write(newValue)
        }()
    }

    // Wait for operations to complete
    time.Sleep(3 * time.Second)
}
