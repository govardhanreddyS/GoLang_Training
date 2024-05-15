package main

import (
	"fmt"
	"time"
)
func main() {
	// Create a channel to communicate events
	eventCh := make(chan string)
	// Start a goroutine to listen for events
	go eventListener(eventCh)
	// Simulate some events being generated
	for i := 0; i < 5; i++ {
		eventCh <- fmt.Sprintf("Event %d", i+1)
		time.Sleep(1 * time.Second) // Simulate some processing time
	}

	// Close the channel to signal that no more events will be sent
	close(eventCh)

	// Wait for a bit to allow the event listener to process remaining events
	time.Sleep(2 * time.Second)
}

func eventListener(ch <-chan string) {
	// Iterate over the events received from the channel
	for event := range ch {
		fmt.Println("Received event:", event)
		processEvent(event)
	}
	fmt.Println("Event listener has finished.")
}

func processEvent(event string) {
	// Simulate processing of the event
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Processed event:", event)
}
