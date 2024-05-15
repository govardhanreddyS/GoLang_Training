package main

import (
	"fmt"
	"time"
)

const (
	numProducers = 3
	numConsumers = 3
	maxJobs      = 5
)

type Job struct {
	ID  int
	Msg string
}

func main() {
	jobQueue := make(chan Job, 10) // Buffered channel to hold jobs

	// Start producers
	for i := 0; i < numProducers; i++ {
		go producer(i, jobQueue)
	}

	// Start consumers
	for i := 0; i < numConsumers; i++ {
		go consumer(i, jobQueue)
	}

	time.Sleep(5 * time.Second) // Let the simulation run for a while
}

func producer(id int, jobs chan<- Job) {
	for i := 0; i < maxJobs; i++ {
		job := Job{ID: i, Msg: fmt.Sprintf("Producer %d: Job %d", id, i)}
		jobs <- job
		time.Sleep(100 * time.Millisecond) // Simulate job creation time
	}
}

func consumer(id int, jobs <-chan Job) {
	for job := range jobs {
		fmt.Printf("Consumer %d: Processing Job %d: %s\n", id, job.ID, job.Msg)
		time.Sleep(200 * time.Millisecond) // Simulate job processing time
	}
}
/*
In this example:
We have a Job struct representing a unit of work.
There are multiple producer goroutines (numProducers) and consumer goroutines (numConsumers) that interact through a shared channel jobQueue.
Producers continuously generate jobs and send them to the jobQueue channel.
Consumers continuously receive jobs from the jobQueue channel and process them.
Each producer and consumer operates concurrently, and they access the shared channel without locks.
We set a limit (maxJobs) on the number of jobs each producer creates to simulate a finite workload.
This example demonstrates a scenario where multiple producers and consumers work concurrently without the risk of starvation. 
Using channels for communication ensures that both producers and consumers have fair access to the shared resource (the jobQueue channel), allowing them to work efficiently without one type of worker dominating access.
*/