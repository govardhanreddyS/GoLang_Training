package main

import (
    "fmt"
    "runtime"
    "sync"
)

// Worker represents a single worker that processes face data
func worker(id int, jobs <-chan interface{}, results chan<- interface{}, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        // Simulate processing
        // Replace this with your actual face data processing logic
        fmt.Printf("Worker %d processing face data: %v\n", id, job)
        // Simulate result
        result := fmt.Sprintf("Processed by worker %d: %v", id, job)
        results <- result
    }
}

func main() {
    // Set GOMAXPROCS to utilize multiple CPUs
    numCPUs := runtime.NumCPU()
    runtime.GOMAXPROCS(numCPUs)
    
    // Define the number of workers
    numWorkers := numCPUs * 2 // Adjust as needed
    
    // Create channels for sending jobs and receiving results
    jobs := make(chan interface{}, numWorkers*2) // Buffer jobs to reduce blocking
    results := make(chan interface{}, numWorkers*2)
    
    // Create a WaitGroup to wait for all workers to finish
    var wg sync.WaitGroup
    
    // Launch workers
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg)
    }
    
    // Generate face data (replace this with your actual data source)
    faceData := []interface{}{"FaceData1", "FaceData2", "FaceData3", "FaceData4", "FaceData5"}
    
    // Send jobs to workers
    go func() {
        defer close(jobs)
        for _, data := range faceData {
            jobs <- data
        }
    }()
    
    // Collect results
    go func() {
        for result := range results {
            fmt.Println(result)
        }
    }()
    
    // Wait for all workers to finish
    wg.Wait()
    
    fmt.Println("All face data processed successfully.")
}
