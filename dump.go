package main

import (
    "sync"
    "testing"
)

// Function to perform task with channels
func taskWithChannel(ch chan int, n int) {
    // Send value through the channel
    ch <- n
}

// Function to perform task with wait group
func taskWithWaitGroup(wg *sync.WaitGroup, n int) {
    defer wg.Done()
    // Perform task
}

// Benchmark using channels
func BenchmarkChannel(b *testing.B) {
    // Create a channel
    ch := make(chan int)

    // Start the benchmark
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Perform task with channel
        go taskWithChannel(ch, i)
        // Receive value from the channel
        <-ch
    }
}

// Benchmark using wait group
func BenchmarkWaitGroup(b *testing.B) {
    // Create a wait group
    var wg sync.WaitGroup

    // Start the benchmark
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Increment wait group counter
        wg.Add(1)
        // Perform task with wait group
        go taskWithWaitGroup(&wg, i)
    }
    // Wait for all tasks to finish
    wg.Wait()
}
