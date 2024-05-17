package main

import (
    "sync"
    "testing"
)

const numWorkers = 100
const numJobs = 100000

func BenchmarkMutexLock(b *testing.B) {
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup
    var mu sync.Mutex

    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := range jobs {
                // Simulate some work
                result := j * 2

                // Use mutex lock to protect shared resource (results)
                mu.Lock()
                results <- result
                mu.Unlock()
            }
        }()
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()

    close(results)
}

func BenchmarkOptimizedChannel(b *testing.B) {
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup

    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case j, ok := <-jobs:
                    if !ok {
                        return
                    }
                    // Simulate some work
                    result := j * 2

                    // Send result to results channel
                    results <- result
                }
            }
        }()
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()

    close(results)
}
