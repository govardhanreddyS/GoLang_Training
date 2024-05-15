package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	numProducers = 3
	numConsumers = 2
	bufferSize   = 5
)

func producer(id int, outCh chan<- int, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case outCh <- rand.Intn(100):
			fmt.Printf("Producer %d: Sent data\n", id)
		case <-stop:
			fmt.Printf("Producer %d: Received stop signal\n", id)
			return
		default:
			fmt.Printf("Producer %d: Buffer full. Waiting...\n", id)
			time.Sleep(time.Second) // Wait for consumer to process data
		}
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate varying processing time
	}
}

func consumer(id int, inCh <-chan int, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case data := <-inCh:
			fmt.Printf("Consumer %d: Received data %d\n", id, data)
		case <-stop:
			fmt.Printf("Consumer %d: Received stop signal\n", id)
			return
		default:
			fmt.Printf("Consumer %d: Buffer empty. Waiting...\n", id)
			time.Sleep(time.Second) // Wait for producer to send data
		}
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate varying processing time
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(numProducers + numConsumers)

	ringBuffer := make(chan int, bufferSize)
	stop := make(chan struct{})

	// Start producers
	for i := 1; i <= numProducers; i++ {
		go producer(i, ringBuffer, stop, &wg)
	}

	// Start consumers
	for i := 1; i <= numConsumers; i++ {
		go consumer(i, ringBuffer, stop, &wg)
	}

	// Listen for SIGINT (Ctrl+C) and SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	<-sig

	// Send stop signal to producers and consumers
	close(stop)

	// Wait for all producers and consumers to finish
	wg.Wait()

	close(ringBuffer)
}

