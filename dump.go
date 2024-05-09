package main

import (
	"fmt"
	"time"
)

func goroutineA(ch chan bool, quit chan struct{}) {
	for {
		select {
		case <-ch:
			fmt.Println("Goroutine A is triggered")
		case <-quit:
			fmt.Println("Goroutine A is stopping")
			return
		}
	}
}

func goroutineB(ch chan bool, quit chan struct{}) {
	for {
		select {
		case <-ch:
			fmt.Println("Goroutine B is triggered")
		case <-quit:
			fmt.Println("Goroutine B is stopping")
			return
		}
	}
}
func main() {
	trigger := make(chan bool)
	quit := make(chan struct{})

	go goroutineA(trigger, quit)
	go goroutineB(trigger, quit)

	// Simulating conditions to trigger different goroutines
	go func() {
		for {
			time.Sleep(2 * time.Second)
			trigger <- true // Trigger goroutine A
			time.Sleep(3 * time.Second)
			trigger <- true // Trigger goroutine B
		}
	}()

	// Stop the goroutines after a certain period
	time.Sleep(15 * time.Second)
	close(quit)

	// Keep the main goroutine alive for a while to observe the output
	time.Sleep(2 * time.Second)
}

