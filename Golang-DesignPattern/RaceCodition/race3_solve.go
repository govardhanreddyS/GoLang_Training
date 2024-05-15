package main

import (
    "fmt"
    "sync"
)

var (
    data  []int
    wg    sync.WaitGroup
    mutex sync.Mutex
)

func main() {
    ch := make(chan int)

    wg.Add(2)
    go sendData(ch)
    go receiveData(ch)

    wg.Wait()

    fmt.Println("Final Data:", data)
}

func sendData(ch chan int) {
    defer wg.Done()
    for i := 0; i < 1000; i++ {
        ch <- i
    }
    close(ch)
}

func receiveData(ch chan int) {
    defer wg.Done()
    for num := range ch {
        mutex.Lock()
        data = append(data, num)
        mutex.Unlock()
    }
}
