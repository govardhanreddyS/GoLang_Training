/*
Mutex/RWMutex
Mutex, is used to protect critical area and shared resource in the scenario of multi routines.

what’s the difference between sync.Mutex and sync.RWMutex?
In general, RWMutex is more detail for read and write level than Mutex.

Mutex: it’s a global mutex lock for a shared resource. No mater the action is read or write.

RWMutex: it has to 2 locks with read lock and write lock. The below list show you how the 2 locks work

When a shared resource is locked by write-lock, then the write-lock and read-lock of other routines 
will be suspended.
When a shared resource is locked by read-lock. Then the write-lock of other routine will be suspended. 
But the read-lock will NOT be suspended of other routine. That means that multiple read-lock action will be executed in concurrency scenarios.
In general, the write-lock has a high priority than read-lock. So RWMutex is very useful for the scenario about writing less but read more.
*/

package main

import (
    "fmt"
    "sync"
)

type ConcurrentMap struct {
    data map[string]int
    mu   sync.RWMutex
}

func NewConcurrentMap() *ConcurrentMap {
    return &ConcurrentMap{
        data: make(map[string]int),
    }
}

func (cm *ConcurrentMap) Set(key string, value int) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.data[key] = value
}

func (cm *ConcurrentMap) Get(key string) (int, bool) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    value, ok := cm.data[key]
    return value, ok
}

func main() {
    cm := NewConcurrentMap()

    var wg sync.WaitGroup

    // Perform concurrent writes
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(key string, value int) {
            defer wg.Done()
            cm.Set(fmt.Sprintf("key%d", value), value)
        }(fmt.Sprintf("key%d", i), i)
    }

    // Perform concurrent reads
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(key string) {
            defer wg.Done()
            if value, ok := cm.Get(key); ok {
                fmt.Printf("Value for %s: %d\n", key, value)
            } else {
                fmt.Printf("Key %s not found\n", key)
            }
        }(fmt.Sprintf("key%d", i))
    }

    // Wait for all operations to complete
    wg.Wait()
}
