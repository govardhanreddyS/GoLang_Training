package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	const runtime = 1 * time.Second

	// Buffered channel to signal which worker should execute next
	workerCh := make(chan int, 1)

	// Pre-fill the channel with a signal for the greedy worker to start
	workerCh <- 1

	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			select {
			case <-workerCh:
				time.Sleep(3 * time.Nanosecond)
				count++
				workerCh <- 2 // Signal to the polite worker
			}
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			select {
			case <-workerCh:
				time.Sleep(1 * time.Nanosecond)
				count++
				workerCh <- 1 // Signal to the greedy worker
			}
		}

		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()

	wg.Wait()
}
