package main

import (
    "fmt"
    "sync"
    "time"
)

// Job represents a task to be performed by a worker
type Job struct {
    ID       int
    TaskType string
    // You can add more fields as needed
}

// Result represents the result of a completed job
type Result struct {
    JobID int
    Data  interface{}
    // You can add more fields as needed
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, cancel <-chan struct{}) {
    defer wg.Done()
    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                return // Jobs channel closed, exit the goroutine
            }
            // Simulate work based on job type
            switch job.TaskType {
            case "task1":
                time.Sleep(time.Second)
                fmt.Printf("Worker %d finished job %d (TaskType: %s)\n", id, job.ID, job.TaskType)
                // Send result
                results <- Result{JobID: job.ID, Data: "Result for task1"}
            case "task2":
                time.Sleep(2 * time.Second)
                fmt.Printf("Worker %d finished job %d (TaskType: %s)\n", id, job.ID, job.TaskType)
                // Send result
                results <- Result{JobID: job.ID, Data: "Result for task2"}
            }
        case <-cancel:
            fmt.Printf("Worker %d canceled\n", id)
            return // Exit the goroutine if canceled
        }
    }
}

func master(numWorkers int, numJobs int) {
    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    cancel := make(chan struct{})
    var wg sync.WaitGroup

    // Start workers
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg, cancel)
    }

    // Send jobs to workers
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j, TaskType: fmt.Sprintf("task%d", (j % 2) + 1)} // Alternate between task1 and task2
    }
    close(jobs)

    // Collect results
    go func() {
        for result := range results {
            fmt.Printf("Result for job %d: %v\n", result.JobID, result.Data)
        }
    }()

    // Wait for all workers to finish
    wg.Wait()
    close(results)
    fmt.Println("All workers finished")
}

func main() {
    numWorkers := 3
    numJobs := 5
    master(numWorkers, numJobs)
}

