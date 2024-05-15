package main

import (
	"fmt"
	"time"
)

// Message represents a message in a messaging system.
type Message struct {
	ID        int
	Sender    string
	Recipient string
	Content   string
	Timestamp time.Time
}

func generator(messages []Message, out chan<- Message) {
	defer close(out)
	for _, msg := range messages {
		out <- msg
	}
}

func consumer(id int, in <-chan Message) {
	for msg := range in {
		fmt.Printf("Consumer %d received message from %s to %s:\n", id, msg.Sender, msg.Recipient)
		fmt.Printf("Content: %s\n", msg.Content)
		fmt.Printf("Timestamp: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println("---------------------------")
		time.Sleep(time.Second) // Simulating some processing
	}
	fmt.Printf("Consumer %d finished\n", id)
}

func main() {
	messages := []Message{
		{ID: 1, Sender: "Alice", Recipient: "Bob", Content: "Hello Bob!", Timestamp: time.Now()},
		{ID: 2, Sender: "Bob", Recipient: "Alice", Content: "Hi Alice!", Timestamp: time.Now()},
		{ID: 3, Sender: "Charlie", Recipient: "David", Content: "Hey David, how are you?", Timestamp: time.Now()},
		{ID: 4, Sender: "David", Recipient: "Charlie", Content: "I'm good, thanks!", Timestamp: time.Now()},
		{ID: 5, Sender: "Eve", Recipient: "Frank", Content: "Meeting at 3 PM today.", Timestamp: time.Now()},
	}
	numConsumers := 3

	// Create channels
	messageCh := make(chan Message)

	// Start generator
	go generator(messages, messageCh)

	// Start consumers
	for i := 0; i < numConsumers; i++ {
		go consumer(i, messageCh)
	}

	// Wait for a key press before exiting
	fmt.Println("Press any key to exit...")
	var input string
	fmt.Scanln(&input)
	fmt.Println("Exiting...")
}
