package main

import (
    "fmt"
    "sync"
)

var (
    numbers = []int{}
    wg      sync.WaitGroup
)

func main() {
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go appendToSlice(i)
    }
    wg.Wait()

    fmt.Println("Final Slice Length:", len(numbers))
}

func appendToSlice(num int) {
    defer wg.Done()
    numbers = append(numbers, num)
}
