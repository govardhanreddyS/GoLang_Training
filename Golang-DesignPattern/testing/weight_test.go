package main

import (
    "sync"
    "testing"
)

const numWorkers = 100000 // Number of workers

func worker(wg *sync.WaitGroup, result chan int) {
    defer wg.Done()
    // Simulate some work
    sum := 0
    for i := 0; i < 1000; i++ {
        sum += i
    }
    // Send result to channel
    result <- sum
}

func BenchmarkWaitGroup(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var wg sync.WaitGroup
        result := make(chan int, numWorkers)

        // Launch workers
        for i := 0; i < numWorkers; i++ {
            wg.Add(1)
            go worker(&wg, result)
        }

        // Wait for all workers to finish
        wg.Wait()
        close(result)
    }
}
