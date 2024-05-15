package main

import (
	"fmt"
	"sync"
)

// Job represents a unit of work
type Job struct {
	ID  int
	Msg string
}

// Worker represents a worker that performs jobs
type Worker struct {
	ID           int
	JobQueue     chan Job
	WorkerWG     *sync.WaitGroup
	WorkerWGDone *sync.WaitGroup
}

// NewWorker creates a new worker with the given ID and assigns it to the job queue
func NewWorker(id int, jobQueue chan Job, workerWG, workerWGDone *sync.WaitGroup) *Worker {
	return &Worker{
		ID:           id,
		JobQueue:     jobQueue,
		WorkerWG:     workerWG,
		WorkerWGDone: workerWGDone,
	}
}

// Start starts the worker to process jobs
func (w *Worker) Start() {
	go func() {
		defer w.WorkerWGDone.Done()
		for job := range w.JobQueue {
			fmt.Printf("Worker %d is processing job %d: %s\n", w.ID, job.ID, job.Msg)
		}
	}()
}

// Pool represents a pool of workers
type Pool struct {
	Workers      []*Worker
	JobQueue     chan Job
	WorkerWG     *sync.WaitGroup
	WorkerWGDone *sync.WaitGroup
}

// NewPool creates a new pool with the given number of workers
func NewPool(numWorkers int) *Pool {
	jobQueue := make(chan Job)
	workerWG := &sync.WaitGroup{}
	workerWGDone := &sync.WaitGroup{}

	pool := &Pool{
		JobQueue:     jobQueue,
		WorkerWG:     workerWG,
		WorkerWGDone: workerWGDone,
	}

	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, jobQueue, workerWG, workerWGDone)
		pool.Workers = append(pool.Workers, worker)
		workerWG.Add(1)
		worker.Start()
	}

	return pool
}

// SubmitJob submits a job to the pool for processing
func (p *Pool) SubmitJob(job Job) {
	p.JobQueue <- job
}

// Close shuts down the pool and waits for all workers to finish processing
func (p *Pool) Close() {
	close(p.JobQueue)
	p.WorkerWG.Wait()
}

func main() {
	pool := NewPool(3)

	// Submitting jobs to the pool
	for i := 1; i <= 5; i++ {
		pool.SubmitJob(Job{ID: i, Msg: fmt.Sprintf("Message %d", i)})
	}

	// Closing the pool and waiting for all workers to finish
	pool.Close()
}
