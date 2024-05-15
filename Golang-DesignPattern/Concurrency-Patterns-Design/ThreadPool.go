package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID int
	// Add any other fields related to the task
}

type Worker struct {
	ID          int
	taskChannel chan Task
	quit        chan bool
	wg          *sync.WaitGroup
}

func NewWorker(ID int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:          ID,
		taskChannel: make(chan Task),
		quit:        make(chan bool),
		wg:          wg,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-w.taskChannel:
				// Simulate task processing
				fmt.Printf("Worker %d: Started task %d\n", w.ID, task.ID)
				time.Sleep(time.Second) // Simulate processing time
				fmt.Printf("Worker %d: Completed task %d\n", w.ID, task.ID)
				w.wg.Done() // Mark task as completed
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Pool struct {
	workers    []*Worker
	currentIdx int
	wg         *sync.WaitGroup
	mutex      *sync.Mutex
}

func NewPool(maxWorkers int) *Pool {
	pool := &Pool{
		wg:    &sync.WaitGroup{},
		mutex: &sync.Mutex{},
	}

	for i := 1; i <= maxWorkers; i++ {
		worker := NewWorker(i, pool.wg)
		pool.workers = append(pool.workers, worker)
		worker.Start()
	}

	return pool
}

func (p *Pool) AddTask(task Task) {
	p.wg.Add(1)
	go func() {
		worker := p.getWorker()
		worker.taskChannel <- task
	}()
}

func (p *Pool) getWorker() *Worker {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	worker := p.workers[p.currentIdx]
	p.currentIdx = (p.currentIdx + 1) % len(p.workers)
	return worker
}

func main() {
	pool := NewPool(5) // Create a pool with 5 workers

	// Example: Add tasks to the pool
	for i := 1; i <= 10; i++ {
		pool.AddTask(Task{ID: i})
	}

	pool.wg.Wait() // Wait for all tasks to be completed

	fmt.Println("All tasks completed.")
}
