package main

import (
    "fmt"
    "time"
)

func worker(id int, jobs <-chan int, done chan<- bool) {
    for job := range jobs {
        // Simulate some work
        time.Sleep(time.Second)
        fmt.Printf("Worker %d finished job %d\n", id, job)
    }
    done <- true // Notify that this worker has finished
}

func master(numWorkers int) {
    jobs := make(chan int, 10)
    done := make(chan bool)

    // Start workers
    for i := 1; i <= numWorkers; i++ {
        go worker(i, jobs, done)
    }

    // Send jobs to workers
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    // Wait for all workers to finish
    for i := 0; i < numWorkers; i++ {
        <-done
    }
    fmt.Println("All workers finished")
}

func main() {
    numWorkers := 3
    master(numWorkers)
}
