package main

import (
    "fmt"
    "sync"
)

var counter = 0
var mutex sync.Mutex

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go incrementCounter(&wg)
    }
    wg.Wait()
    fmt.Println("Final Counter Value:", counter)
}

func incrementCounter(wg *sync.WaitGroup) {
    mutex.Lock()
    counter = counter + 1
    mutex.Unlock()
    wg.Done()
}
