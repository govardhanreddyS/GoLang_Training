package main
/*
Introduced a Task struct representing a unit of work to be performed by a goroutine.
Added a Perform method to the Task struct to simulate the execution of the task.
Increased clarity by separating the task generation and execution logic into distinct parts.
Improved flexibility by allowing the number of goroutines and the number of tasks to be easily configurable.
Ensured that all tasks are sent to the goroutines before closing the task channel to avoid premature closure.
Improved comments to provide better understanding of the code.
*/
import (
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work to be performed by a goroutine.
type Task struct {
	ID int
}

// Perform simulates the execution of a task.
func (t Task) Perform() {
	fmt.Printf("Task %d is being executed\n", t.ID)
	time.Sleep(1 * time.Second) // Simulating some time-consuming operation
	fmt.Printf("Task %d has completed\n", t.ID)
}

func main() {
	// Create a waitgroup to track the completion of all goroutines.
	var wg sync.WaitGroup

	// Define the number of goroutines to run concurrently.
	numGoroutines := 10

	// Create a channel to limit the number of goroutines that can be running at the same time.
	semaphore := make(chan struct{}, numGoroutines)

	// Create a channel to send tasks to goroutines.
	taskChannel := make(chan Task)

	// Define a function to be executed by each goroutine.
	goroutineFunc := func(i int) {
		defer wg.Done()

		// Receive tasks from the task channel until it is closed.
		for task := range taskChannel {
			// Acquire a semaphore token.
			semaphore <- struct{}{}

			// Perform the task.
			task.Perform()

			// Release the semaphore token.
			<-semaphore
		}
	}

	// Start goroutines.
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go goroutineFunc(i)
	}

	// Generate tasks and send them to the task channel.
	go func() {
		for i := 0; i < 20; i++ {
			taskChannel <- Task{ID: i}
		}
		close(taskChannel) // Close the task channel when all tasks are sent.
	}()

	// Wait for all goroutines to finish.
	wg.Wait()

	fmt.Println("All tasks have completed.")
}

