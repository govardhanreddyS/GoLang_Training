package main

import (
	"fmt"
	"runtime"
	//"sync"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID  int
	Job string
}

// Worker represents a worker that executes tasks
type Worker struct {
	ID         int
	taskChan   chan Task
	quitSignal chan struct{}
}

// NewWorker creates a new worker
func NewWorker(id int, taskChan chan Task) *Worker {
	return &Worker{
		ID:         id,
		taskChan:   taskChan,
		quitSignal: make(chan struct{}),
	}
}

// Start starts the worker
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-w.taskChan:
				fmt.Printf("Worker %d is processing task %d: %s\n", w.ID, task.ID, task.Job)
			case <-w.quitSignal:
				return
			}
		}
	}()
}

// Stop stops the worker
func (w *Worker) Stop() {
	close(w.quitSignal)
}

// Dispatcher represents a dispatcher that distributes tasks to workers
type Dispatcher struct {
	workers []*Worker
}

// NewDispatcher creates a new dispatcher with the specified number of workers
func NewDispatcher(numWorkers int, taskChan chan Task) *Dispatcher {
	dispatcher := &Dispatcher{}
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, taskChan)
		dispatcher.workers = append(dispatcher.workers, worker)
	}
	return dispatcher
}

// Start starts all the workers in the dispatcher
func (d *Dispatcher) Start() {
	for _, worker := range d.workers {
		worker.Start()
	}
}

// Stop stops all the workers in the dispatcher
func (d *Dispatcher) Stop() {
	for _, worker := range d.workers {
		worker.Stop()
	}
}

// SubmitTask submits a task to be executed by one of the workers
func (d *Dispatcher) SubmitTask(task Task) {
	// Choose a worker to execute the task (simple round-robin scheduling)
	worker := d.workers[task.ID%len(d.workers)]
	worker.taskChan <- task
}

func main() {
	// Get the number of CPU cores
	numCores := runtime.NumCPU()

	// Use twice the number of CPU cores as the number of workers
	numWorkers := 4* numCores

	// Create a channel for tasks
	taskChan := make(chan Task)

	// Create a dispatcher with the specified number of workers
	dispatcher := NewDispatcher(numWorkers, taskChan)
	startTime := time.Now()
	// Start the dispatcher
	dispatcher.Start()

	// Submit some tasks
	numTasks := 20000
	for i := 1; i <= numTasks; i++ {
		task := Task{ID: i, Job: fmt.Sprintf("Message %d", i)}
		dispatcher.SubmitTask(task)
	}

	// Wait for all tasks to be processed
	//time.Sleep(1 * time.Second)

	// Stop the dispatcher
	dispatcher.Stop()

	// Calculate and print the time taken
	endTime := time.Now()
	fmt.Printf("Time taken: %s\n", endTime.Sub(startTime))
}
