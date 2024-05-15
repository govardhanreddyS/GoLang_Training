package main

import (
	"fmt"
	"sync"
	//"time"
)

// orDone ensures that the returned channel only receives values from c until either done is closed or c is closed
func orDone(done <-chan struct{}, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return valStream
}

// bridge orchestrates multiple input channels and sends their values to a single output channel
func bridge(done <-chan struct{}, chanStreams ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	output := make(chan interface{})

	outputStream := func(c <-chan interface{}) {
		defer wg.Done()
		for val := range orDone(done, c) {
			select {
			case output <- val:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(chanStreams))
	for _, c := range chanStreams {
		go outputStream(c)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// genStreams generates a sequence of channels, each containing a varying number of integer values
func genStreams(n int) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for i := 0; i < n; i++ {
			stream <- i+1
		}
	}()
	return stream
}

func main() {
	// Create a cancellation channel
	cancel := make(chan struct{})
	defer close(cancel) // Ensure cancellation channel is closed when main exits

	// Use bridge to receive values from multiple streams and print them
	output := bridge(cancel, genStreams(3), genStreams(5), genStreams(2))

	// Print received values
	for value := range output {
		fmt.Printf("%v ", value)
	}

	fmt.Println("\nProcessing complete")
}
