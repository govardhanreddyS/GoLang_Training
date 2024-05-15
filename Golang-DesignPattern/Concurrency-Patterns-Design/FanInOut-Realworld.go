package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Generator generates random numbers and sends them to a channel.
func Generator(out chan<- int) {
	defer close(out)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		// Generate a random number between 1 and 100
		num := rand.Intn(100) + 1
		out <- num
	}
}

// Worker processes numbers received from an input channel and sends the result to an output channel.
func Worker(id int, in <-chan int, out chan<- int) {
	defer close(out)
	for num := range in {
		// Perform some processing (e.g., doubling the number)
		result := num * 2
		fmt.Printf("Worker %d processed %d and produced %d\n", id, num, result)
		out <- result
	}
}

// Merge merges results from multiple workers into a single channel.
func Merge(workers []<-chan int) <-chan int {
	mergeOut := make(chan int)
	go func() {
		defer close(mergeOut)
		for _, worker := range workers {
			for num := range worker {
				mergeOut <- num
			}
		}
	}()
	return mergeOut
}

func main() {
	// Create channels
	generatorOut := make(chan int)
	var workerOuts []chan int

	// Start generator
	go Generator(generatorOut)

	// Start workers
	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		workerOut := make(chan int)
		workerOuts = append(workerOuts, workerOut)
		go Worker(i, generatorOut, workerOut)
	}

	// Collect results from workers
	var workersResult []<-chan int
	for _, workerOut := range workerOuts {
		workersResult = append(workersResult, workerOut)
	}

	// Merge results from workers
	mergeOut := Merge(workersResult)

	// Process merged results
	for result := range mergeOut {
		fmt.Printf("Result: %d\n", result)
	}

	fmt.Println("Done")
}
/*

The above code can be used in various scenarios where there's a need to process data concurrently with multiple workers and merge the results. Here's a hypothetical context where the code could be applied:

Scenario: Data Processing Pipeline

Imagine you're developing a data processing pipeline for a system that needs to analyze and process incoming data streams. The system receives streams of numeric data from various sources and needs to perform some computation on each data point, such as transforming, filtering, or aggregating the data. The processed data points are then used for further analysis or stored in a database for reporting purposes.

Context:

Data Generation: The Generator function simulates the initial data source, generating random numbers that represent incoming data points.

Data Processing Workers: The Worker function represents the individual processing steps in the pipeline. Each worker performs a specific computation on the data points received from the Generator. For example, one worker might double the value of each data point, while another worker might calculate the square of each data point.

Concurrency: By using goroutines, the code enables concurrent processing of data points with multiple workers. Each worker operates independently, processing data points concurrently without waiting for others to finish.

Merge Results: The Merge function combines the results from all workers into a single channel. This enables the system to aggregate the processed data points from different processing steps into a unified stream for further analysis or storage.

Usage:

The code can be integrated into the data processing pipeline of the system to handle the processing of incoming data streams efficiently. It allows the system to scale horizontally by adding more processing workers as needed to handle increased data volumes.

For example, in a real-world application:

The Generator function could be replaced with a component that reads data from sensors, network streams, or database queries.

The Worker function could represent various data processing tasks, such as data transformation, filtering, feature extraction, or anomaly detection.

The Merge function could feed the processed data into subsequent stages of the pipeline, such as machine learning models, statistical analysis modules, or database storage components.

Overall, this code provides a flexible and scalable solution for processing data streams in real-time, enabling the system to efficiently handle large volumes of data while maintaining high throughput and low latency.
*/