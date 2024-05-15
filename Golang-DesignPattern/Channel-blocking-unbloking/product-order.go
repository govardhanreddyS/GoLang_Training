package main

import (
	"fmt"
	"time"
)

// simulateOrder represents an order in the online store.
type simulateOrder struct {
	productType string
}

// goroutineA simulates processing orders for physical products.
func goroutineA(ch chan simulateOrder, quit chan struct{}) {
	for {
		select {
		case order := <-ch:
			fmt.Printf("Processing order for physical product: %v\n", order)
			// Simulate processing time
			time.Sleep(2 * time.Second)
			fmt.Printf("Order processed for physical product: %v\n", order)
		case <-quit:
			fmt.Println("Goroutine A is stopping")
			return
		}
	}
}

// goroutineB simulates processing orders for digital products.
func goroutineB(ch chan simulateOrder, quit chan struct{}) {
	for {
		select {
		case order := <-ch:
			fmt.Printf("Processing order for digital product: %v\n", order)
			// Simulate processing time
			time.Sleep(1 * time.Second)
			fmt.Printf("Order processed for digital product: %v\n", order)
		case <-quit:
			fmt.Println("Goroutine B is stopping")
			return
		}
	}
}

func main() {
	orderChannel := make(chan simulateOrder)
	quit := make(chan struct{})

	go goroutineA(orderChannel, quit)
	go goroutineB(orderChannel, quit)

	// Simulate incoming orders
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(2 * time.Second)
			orderChannel <- simulateOrder{"Physical Product " + fmt.Sprint(i)}
		}
		for i := 1; i <= 5; i++ {
			time.Sleep(3 * time.Second)
			orderChannel <- simulateOrder{"Digital Product " + fmt.Sprint(i)}
		}
	}()

	// Stop the goroutines after processing orders
	time.Sleep(20 * time.Second)
	close(quit)

	// Keep the main goroutine alive for a while to observe the output
	time.Sleep(2 * time.Second)
}
