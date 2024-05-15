package main

import (
    "testing"
)

const cnumWorkers = 100000 // Number of workers

func cworker(result chan int) {
    // Simulate some work
    sum := 0
    for i := 0; i < 1000; i++ {
        sum += i
    }
    // Send result to channel
    result <- sum
}

func BenchmarkChannels(b *testing.B) {
    for i := 0; i < b.N; i++ {
        result := make(chan int, cnumWorkers)

        // Launch workers
        for i := 0; i < cnumWorkers; i++ {
            go cworker(result)
        }

        // Collect results
        total := 0
        for i := 0; i < cnumWorkers; i++ {
            total += <-result
        }
        close(result)
    }
}
