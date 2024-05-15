package main
/*
In this example:
We define a worker function that simulates processing tasks from a task queue. Each worker runs in its goroutine and listens for tasks from the taskQueue channel. It also listens for cancellation signals from the ctx.
In the main function, we create a context ctx and a cancellation function cancel using context.WithCancel. This context is passed to all workers, allowing us to cancel their execution when needed.
We start multiple workers (controlled by the numWorkers variable) in separate goroutines.
We add tasks to the taskQueue channel.
After some simulated processing time, we close the taskQueue channel to indicate that no more tasks will be added.
We wait for all workers to finish their tasks using wg.Wait().
This example demonstrates how to manage concurrency using context for graceful shutdown and cancellation handling in a moderately complex scenario with multiple workers and task processing.
*/
import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, taskQueue <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d: Starting\n", id)

	for {
		select {
		case task, ok := <-taskQueue:
			if !ok {
				fmt.Printf("Worker %d: Task queue closed. Exiting.\n", id)
				return
			}
			fmt.Printf("Worker %d: Processing task %d\n", id, task)
			time.Sleep(1 * time.Second) // Simulate task processing time
		case <-ctx.Done():
			fmt.Printf("Worker %d: Context canceled. Exiting.\n", id)
			return
		}
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskQueue := make(chan int, 10)
	numWorkers := 3
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, taskQueue, &wg)
	}

	// Add tasks to the queue
	for i := 1; i <= 10; i++ {
		taskQueue <- i
	}

	// Simulate some processing time
	time.Sleep(3 * time.Second)

	// Close the task queue and wait for workers to finish
	close(taskQueue)
	wg.Wait()

	fmt.Println("All workers have finished their tasks.")
}
