package main
/*
We need this design pattern when we required to join asynchronous and synchronus pattern togather 
In this example:
The UserFetcher component simulates fetching user data synchronously from a remote API. It holds the fetched user data in a map.
The DataProcessor component processes the fetched user data asynchronously. It calculates statistics on the user data, such as the length of each user's name.
We create instances of UserFetcher and DataProcessor and use them to fetch and process user data.
The processing of user data is performed asynchronously, allowing the main program to continue executing other tasks concurrently.
After retrieving the processed data asynchronously, we print the statistics.
You can run this code to see how the synchronous fetching and asynchronous processing of data work together in a practical scenario.
*/
import (
	"fmt"
	"sync"
	"time"
)

// Synchronous component for fetching user data
type UserFetcher struct {
	mu   sync.Mutex
	data map[int]string
}

// Asynchronous component for processing user data
type DataProcessor struct {
	ch chan map[int]int
}

// NewUserFetcher creates a new UserFetcher component
func NewUserFetcher() *UserFetcher {
	return &UserFetcher{
		data: make(map[int]string),
	}
}

// NewDataProcessor creates a new DataProcessor component
func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		ch: make(chan map[int]int),
	}
}

// FetchUserData simulates fetching user data synchronously
func (uf *UserFetcher) FetchUserData() {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	// Simulate fetching user data from a remote API
	time.Sleep(2 * time.Second)

	// Populate user data map
	uf.data[1] = "Alice"
	uf.data[2] = "Bob"
	uf.data[3] = "Charlie"
}

// ProcessData asynchronously processes the fetched user data
func (dp *DataProcessor) ProcessData(uf *UserFetcher) {
	go func() {
		// Retrieve user data from UserFetcher
		uf.mu.Lock()
		data := uf.data
		uf.mu.Unlock()

		// Simulate processing time
		time.Sleep(3 * time.Second)

		// Perform some processing (e.g., calculating statistics)
		stats := make(map[int]int)
		for id, name := range data {
			stats[id] = len(name)
		}

		// Send processed data through channel
		dp.ch <- stats
	}()
}

func main() {
	// Create synchronous and asynchronous components
	userFetcher := NewUserFetcher()
	dataProcessor := NewDataProcessor()

	// Fetch user data synchronously
	fmt.Println("Fetching user data synchronously...")
	userFetcher.FetchUserData()

	// Process user data asynchronously
	fmt.Println("Processing user data asynchronously...")
	dataProcessor.ProcessData(userFetcher)

	// Retrieve processed data asynchronously
	stats := <-dataProcessor.ch

	// Print processed data
	fmt.Println("Processed data:")
	for id, length := range stats {
		fmt.Printf("User ID: %d, Name Length: %d\n", id, length)
	}

	// Wait for processing to finish
	time.Sleep(4 * time.Second)
}
