// Worker pool benefits:
// - Efficiency because it distributes the work across threads.
// - Flow control: Limit work in flight

// Disadvantage of worker:
// Lifetimes complexity: clean up and idle worker

// Principles:
// Start goroutines whenever you have the concurrent work to do.
// The goroutine should exit as soon as posible the work is done. This helps us
// to clean up the resources and manage the lifetimes correctly.
package main

import (
	"fmt"
	//"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Millisecond)
		fmt.Println("worker", id, "fnished job", j)
		results <- j * 2
	}
}

func workerEfficient(id int, jobs <-chan int, results chan<- int) {
	// sync.WaitGroup helps us to manage the job
	//var wg sync.WaitGroup
	for j := range jobs {

	//	wg.Add(1)
		// we start a goroutine to run the job
		go func(job int) {
			// start the job
			fmt.Println("Efficient worker", id, "started job", job)
			time.Sleep(time.Millisecond*10)
			fmt.Println("Efficient worker", id, "fnished job", job)
			results <- job * 3
		//	wg.Done()

		}(j)

	}
	// With a help to manage the lifetimes of goroutines
	// we can add more handler when a goroutine finished
//	wg.Wait()
}
func main() {
	const numbJobs = 500
	jobs := make(chan int, numbJobs)
	results := make(chan int, numbJobs)

	// 1. Start the worker
	// it is a fixed pool of goroutines receive and perform tasks from a channel

	// In this example, we define a fixed 3 workers
	// they receive the `jobs` from the channel jobs
	// we also naming the worker name with `w` variable.
	for w := 1; w <= 50; w++ {
		
		go worker(w, jobs, results)
		go workerEfficient(w, jobs, results)
	}

	// 2. send the work
	// other goroutine sends the work to the channels

	// in this example, the `main` goroutine sends the work to the channel `jobs`
	for j := 1; j <= numbJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	time.Sleep(time.Millisecond*1000)
	fmt.Println("Closed job")
	for a := 1; a <= numbJobs; a++ {
		<-results
	}
	close(results)

}