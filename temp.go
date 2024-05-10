package main
/*
Worker Pool with Structured Task Data:
In this example, we'll create a worker pool to process structured task data
 concurrently. Imagine a scenario where you have a large set of tasks to 
 execute, and each task is represented by a structured data type.
*/
import (
    "fmt"
    "sync"
)

// Task represents a structured task data
type Task struct {
    ID   int
    Data string
}

func main() {
    tasks := make(chan Task, 100) // Buffered channel for tasks
    results := make(chan Task, 100)
    var wg sync.WaitGroup

    // Worker function to process tasks
    worker := func(id int, tasks <-chan Task, results chan<- Task) {
        defer wg.Done()
        for task := range tasks {
            // Process the task
            task.Data = fmt.Sprintf("Processed: %s", task.Data)
            results <- task
        }
    }

    // Spawn worker goroutines
    numWorkers := 5
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, tasks, results)
    }

    // Generate tasks
    for i := 1; i <= 20; i++ {
        tasks <- Task{ID: i, Data: fmt.Sprintf("Task %d", i)}
    }
    close(tasks)

    // Collect results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Print results
    for result := range results {
        fmt.Println(result.Data)
    }
}
