package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	counter := 0
	//done:=make(chan bool)

	// Goroutine to increment the counter
	increment := func() {
		defer wg.Done()
		for {
			mu.Lock()
			counter++
			fmt.Println("Incremented counter to:", counter)
			cond.Signal() // Signal that the counter has been incremented
			mu.Unlock()
		//	done<-true
			//time.Sleep(time.Second*1)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Sleep for a random duration up to 1 second
		}
	}

	// Goroutine to decrement the counter
	decrement := func() {
		defer wg.Done()
		for {
			//<-done
			mu.Lock()
			for counter == 0 {
				cond.Wait() // Wait until counter is non-zero
			}
			counter--
			fmt.Println("Decremented counter to:", counter)
			mu.Unlock()
			//time.Sleep(time.Second * 3)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		
		}
	}

	// Start the goroutines
	wg.Add(2)
	
	go increment()
	go decrement()

	// Wait for the goroutines to finish
	wg.Wait()
	//close(done)
}
