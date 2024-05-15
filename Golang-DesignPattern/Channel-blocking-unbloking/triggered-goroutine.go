package main

import (
	"fmt"
	"time"
)

func goroutineA(ch chan bool ) {
	for {
		select {
		case <-ch:
			fmt.Println("Goroutine A is triggered")
		
		}
	}
}

func goroutineB(ch chan bool) {
	for {
		select {
		case <-ch:
			fmt.Println("Goroutine B is triggered")
		}
	}
}
func main() {
	trigger := make(chan bool)
	//done := make(chan bool)

	go goroutineA(trigger )
	go goroutineB(trigger )
	// Simulating conditions to trigger different goroutines
	go func() {
		for {
			time.Sleep(2 * time.Second)
			trigger <- true // Trigger goroutine A
			time.Sleep(3 * time.Second)
			trigger <- true // Trigger goroutine B
		}
	}()

	// Keep the main goroutine alive
	select {}
}
