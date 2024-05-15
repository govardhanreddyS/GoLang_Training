package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numProducers = 3
	numConsumers = 2
	bufferSize   = 5
)

func producer(id int, outCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case outCh <- rand.Intn(100):
			fmt.Printf("Producer %d: Sent data\n", id)
		default:
			fmt.Printf("Producer %d: Buffer full. Waiting...\n", id)
			time.Sleep(time.Second) // Wait for consumer to process data
		}
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate varying processing time
	}
}

func consumer(id int, inCh <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range inCh {
		fmt.Printf("Consumer %d: Received data %d\n", id, data)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate processing time
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(numProducers + numConsumers)

	ringBuffer := make(chan int, bufferSize)

	// Start producers
	for i := 1; i <= numProducers; i++ {
		go producer(i, ringBuffer, &wg)
	}

	// Start consumers
	for i := 1; i <= numConsumers; i++ {
		go consumer(i, ringBuffer, &wg)
	}

	// Wait for all producers and consumers to finish
	wg.Wait()

	close(ringBuffer)
}

/*
Explanation:
We define numProducers and numConsumers constants to control the number of producers and consumers.
The producer function generates random data and sends it to the channel outCh. 
If the buffer is full, it waits for some time before retrying.
The consumer function receives data from the channel inCh and processes it.
 It simulates processing time before moving on to the next data.
In the main function, we create a buffered channel ringBuffer.
We start multiple goroutines for producers and consumers, each identified by a unique ID.
The sync.WaitGroup is used to wait for all producers and consumers to finish before closing the channel.
Finally, we close the channel to signal that no more data will be sent.
This example demonstrates how to use a ring buffer in a concurrent setting with multiple producers and consumers, providing a flexible and efficient way to manage data flow.
*/