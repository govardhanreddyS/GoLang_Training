/*
Here's how the example serves its purpose:

Demonstrating Concurrency: 

By using goroutines, we simulate concurrent processing of social media 
posts by multiple workers. Each worker runs independently and processes messages concurrently, allowing for efficient 
use of resources and improved performance.

Decoupling of Components: 
The fanout pattern allows us to 
decouple the generation of messages from their processing.
 Messages are generated and sent to a central channel, from which multiple worker goroutines consume them independently. This decoupling improves modularity and flexibility, as we can easily add or remove worker components without affecting the message generation process.

Scalability: 
The example demonstrates how the system can scale horizontally by adding more instances of worker components. We can increase the number of sentiment analysis workers, topic extraction workers, and storage workers to handle higher loads and improve processing speed.

Parallel Processing: 
Each worker component performs a specific task (sentiment analysis, topic extraction, storage) on incoming messages in parallel. This parallel processing improves throughput and reduces processing time, allowing the system to handle a large volume of messages efficiently.

Simulation of Real-World Scenario: 
While the example simplifies the actual complexity of a social media analytics platform, it captures the essential aspects of processing incoming messages, performing various analyses on them, and storing the results for further analysis. This simulation helps illustrate how the fanout pattern can be applied in a real-world scenario to build scalable and efficient systems.

*/
package main

import (
	"fmt"
	"time"
)

// Message represents a social media post.
type Message struct {
	ID        int
	Platform  string
	Username  string
	Content   string
	Timestamp time.Time
}

// SentimentAnalysisWorker simulates a worker for sentiment analysis.
func SentimentAnalysisWorker(id int, in <-chan Message, out chan<- Message) {
	for msg := range in {
		// Simulate sentiment analysis (just printing for demonstration)
		fmt.Printf("SentimentAnalysisWorker %d analyzing sentiment of message %d\n", id, msg.ID)
		// Simulate processing time
		time.Sleep(time.Second)
		// Send the message to the next worker
		out <- msg
	}
}

// TopicExtractionWorker simulates a worker for topic extraction.
func TopicExtractionWorker(id int, in <-chan Message, out chan<- Message) {
	for msg := range in {
		// Simulate topic extraction (just printing for demonstration)
		fmt.Printf("TopicExtractionWorker %d extracting topics from message %d\n", id, msg.ID)
		// Simulate processing time
		time.Sleep(time.Second)
		// Send the message to the next worker
		out <- msg
	}
}

// StorageWorker simulates a worker for storing processed messages.
func StorageWorker(id int, in <-chan Message) {
	for msg := range in {
		// Simulate storing the message in a database (just printing for demonstration)
		fmt.Printf("StorageWorker %d storing message %d in database\n", id, msg.ID)
		// Simulate database storage time
		time.Sleep(time.Second)
	}
}

func main() {
	// Generate some sample messages
	messages := []Message{
		{ID: 1, Platform: "Twitter", Username: "user1", Content: "Hello world!", Timestamp: time.Now()},
		{ID: 2, Platform: "Facebook", Username: "user2", Content: "Good morning!", Timestamp: time.Now()},
		{ID: 3, Platform: "Instagram", Username: "user3", Content: "Happy Friday!", Timestamp: time.Now()},
		// Add more sample messages as needed
	}

	// Create channels
	messageCh := make(chan Message)
	sentimentCh := make(chan Message)
	topicCh := make(chan Message)

	// Start sentiment analysis workers
	for i := 1; i <= 3; i++ {
		go SentimentAnalysisWorker(i, messageCh, sentimentCh)
	}

	// Start topic extraction workers
	for i := 1; i <= 2; i++ {
		go TopicExtractionWorker(i, sentimentCh, topicCh)
	}

	// Start storage workers
	for i := 1; i <= 2; i++ {
		go StorageWorker(i, topicCh)
	}

	// Simulate generating messages
	go func() {
		for _, msg := range messages {
			messageCh <- msg
			// Simulate message generation time
			time.Sleep(500 * time.Millisecond)
		}
		close(messageCh)
	}()

	// Wait for a key press before exiting
	fmt.Println("Press any key to exit...")
	var input string
	fmt.Scanln(&input)
	fmt.Println("Exiting...")
}
