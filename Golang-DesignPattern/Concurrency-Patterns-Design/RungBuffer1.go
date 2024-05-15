package main

import (
	"fmt"
	"time"
)

func producer(outCh chan<- int) {
	for i := 0; i < 20; i++ {
		fmt.Println("Producer sending:", i)
		select {
		case outCh <- i:
			fmt.Println("Sent to buffer:", i)
		default:
			fmt.Println("Buffer full. Waiting for consumer to process data...")
			time.Sleep(time.Second) // Wait for consumer to process data
			i-- // Retry sending the current value
		}
		time.Sleep(time.Second) // Simulate some processing time
	}
	close(outCh)
}

func consumer(inCh <-chan int) {
	for val := range inCh {
		fmt.Println("Consumer received:", val)
		time.Sleep(2 * time.Second) // Simulate some processing time
	}
}

func main() {
	bufferSize := 5
	ringBuffer := make(chan int, bufferSize)

	go producer(ringBuffer)
	consumer(ringBuffer)
}
