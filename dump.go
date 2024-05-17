package main

import (
    "sync"
    "testing"
)

// Define a task function that simulates some computational work.
func doWork() {
    // Simulate some work by looping for a large number of iterations.
    for i := 0; i < 1000000; i++ {
        // Do some computation.
    }
}

// BenchmarkChannelConcurrency benchmarks the performance of using channels for concurrency.
func BenchmarkChannelConcurrency(b *testing.B) {
    // Start the benchmarking loop.
    for i := 0; i < b.N; i++ {
        // Create a channel to coordinate goroutines.
        ch := make(chan struct{})

        // Spawn multiple goroutines to perform the task concurrently.
        for j := 0; j < 10; j++ {
            go func() {
                // Perform the task.
                doWork()
                // Signal completion to the channel.
                ch <- struct{}{}
            }()
        }

        // Wait for all goroutines to finish by receiving from the channel.
        for j := 0; j < 10; j++ {
            <-ch
        }
    }
}

// BenchmarkWaitGroupConcurrency benchmarks the performance of using wait groups for concurrency.
func BenchmarkWaitGroupConcurrency(b *testing.B) {
    // Start the benchmarking loop.
    for i := 0; i < b.N; i++ {
        // Create a wait group to wait for goroutines to finish.
        var wg sync.WaitGroup

        // Add the number of goroutines to the wait group.
        wg.Add(10)

        // Spawn multiple goroutines to perform the task concurrently.
        for j := 0; j < 10; j++ {
            go func() {
                // Perform the task.
                doWork()
                // Notify the wait group that this goroutine has finished.
                wg.Done()
            }()
        }

        // Wait for all goroutines to finish by calling Wait on the wait group.
        wg.Wait()
    }
}
