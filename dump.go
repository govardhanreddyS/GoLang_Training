package main

import (
	"fmt"
	"sync"
	"time"
)

// Event represents a generic event type
type Event struct {
	Type    string
	Payload interface{}
}

// EventHandler represents a function that handles events
type EventHandler func(event Event)

func main() {
	// Create a channel to communicate events
	eventCh := make(chan Event)

	// Start the event listener
	var wg sync.WaitGroup
	wg.Add(1)
	go eventListener(eventCh, &wg)

	// Subscribe to specific event types
	subscribe(eventCh, "EventTypeA", handleEventTypeA)
	subscribe(eventCh, "EventTypeB", handleEventTypeB)
	subscribe(eventCh, "EventType1", handleEventType1)

	// Simulate generating some events
	for i := 0; i < 20; i++ {
		event := Event{
			Type:    fmt.Sprintf("EventType%d", i+1),
			Payload: fmt.Sprintf("Payload %d", i+1),
		}
		eventCh <- event
		time.Sleep(1 * time.Second) // Simulate some processing time
	}

	// Wait for a bit and then gracefully shutdown the event listener
	time.Sleep(2 * time.Second)
	close(eventCh) // Signal that no more events will be sent
	wg.Wait()     // Wait for the event listener to finish
	fmt.Println("Main goroutine has finished.")
}

// eventListener listens for events from the channel and dispatches them to handlers
func eventListener(ch <-chan Event, wg *sync.WaitGroup) {
	defer wg.Done()

	// Iterate over the events received from the channel
	for event := range ch {
		fmt.Printf("Received event: Type=%s, Payload=%v\n", event.Type, event.Payload)
		dispatchEvent(event)
	}
	fmt.Println("Event listener has finished.")
}

// dispatchEvent dispatches the event to the appropriate handler based on its type
func dispatchEvent(event Event) {
	switch event.Type {
	case "EventTypeA":
		handleEventTypeA(event)
	case "EventTypeB":
		handleEventTypeB(event)
	case  "EventType1" :
		handleEventType1(event)
	default:
		fmt.Println("Unknown event type:", event.Type)
	}
}

// subscribe adds a handler function to the event listener for a specific event type
func subscribe(ch <-chan Event, eventType string, handler EventHandler) {
	go func() {
		for event := range ch {
			if event.Type == eventType {
				handler(event)
			}
		}
	}()
}

// handleEventTypeA is an example handler for EventTypeA events
func handleEventTypeA(event Event) {
	// Simulate processing of EventTypeA event
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Processed EventTypeA event:", event.Payload)
}


func handleEventType1(event Event) {
	// Simulate processing of EventTypeA event
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Processed EventType1 event:", event.Payload)
}
// handleEventTypeB is an example handler for EventTypeB events
func handleEventTypeB(event Event) {
	// Simulate processing of EventTypeB event
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Processed EventTypeB event:", event.Payload)
}

/*
The code you've provided appears to be an event-driven program written in Go. It creates events of different types and payloads, subscribes handlers to specific event types, and processes the events asynchronously.

Here's a breakdown of the code:

Event Struct: Defines a generic event type with Type and Payload fields.

EventHandler Type: Represents a function signature for event handlers.

Main Function:

Creates a channel eventCh to communicate events.
Starts the event listener goroutine with eventListener.
Subscribes event handlers for specific event types.
Generates events of different types and sends them to the channel.
Gracefully shuts down the event listener after all events have been processed.
eventListener Function: Listens for events from the channel and dispatches them to handlers.

Uses a sync.WaitGroup to ensure that the main goroutine waits for the event listener to finish.
Iterates over events received from the channel and calls dispatchEvent to handle them.
dispatchEvent Function: Dispatches events to appropriate handlers based on their types.

Matches event types and calls corresponding handler functions.
subscribe Function: Adds handler functions to the event listener for specific event types.

Creates a goroutine for each event type to listen for events and call the corresponding handler.
Handler Functions:

handleEventTypeA: Simulates processing of EventTypeA events with a 500ms delay.
handleEventType1: Simulates processing of EventType1 events with a 500ms delay.
handleEventTypeB: Simulates processing of EventTypeB events with a 1000ms delay.
Overall, this code demonstrates an event-driven architecture in Go using channels and goroutines for concurrency.
 It's a scalable approach to handle asynchronous events in a modular and decoupled manner. 
Let me know if you have any questions or if there's anything specific you'd like to discuss!
*/
