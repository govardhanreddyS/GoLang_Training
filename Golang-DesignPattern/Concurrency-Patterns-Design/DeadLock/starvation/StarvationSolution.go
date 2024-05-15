package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	const runtime = 1 * time.Second

	// Use a mutex to protect access to the worker turn
	var mu sync.Mutex
	var worker1Turn bool = true

	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			mu.Lock()
			if worker1Turn {
				time.Sleep(3 * time.Nanosecond)
				count++
				worker1Turn = false
			}
			mu.Unlock()
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			mu.Lock()
			if !worker1Turn {
				time.Sleep(1 * time.Nanosecond)
				count++
				worker1Turn = true
			}
			mu.Unlock()
		}

		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()

	wg.Wait()
}
/*
The key change that made the synchronization work was the introduction of
 a mutual exclusion mechanism to control access to the shared variable worker1Turn. 
 This ensures that only one worker can check or update the worker1Turn variable at a time,
 preventing both workers from executing simultaneously and avoiding deadlocks.

Here's a summary of the changes made to the code:

Introduction of Mutex: We used a sync.Mutex named mu to protect access to the 
worker1Turn variable.
Before reading or updating worker1Turn, each worker acquires the lock on mu using 
mu.Lock() and releases it using mu.Unlock() afterward. This ensures that only one worker
can  access the worker1Turn variable at a time.
Conditional Execution: Each worker checks the value of worker1Turn within a critical section 
protected by the mutex. If it's their turn to execute (worker1Turn is true for the polite worker 
and false for the greedy worker), they execute their task and update worker1Turn to allow 
the other worker to execute next.
Balanced Workload: By controlling access to the worker1Turn variable using the mutex, we 
ensure that the execution of both workers alternates. This balanced approach allows both
 workers to execute an almost equal number of work loops, preventing one worker from 
 monopolizing the shared resource.
Overall, these changes ensure fair and synchronized execution between the two workers, 
leading to a successful resolution of the deadlock issu
*/