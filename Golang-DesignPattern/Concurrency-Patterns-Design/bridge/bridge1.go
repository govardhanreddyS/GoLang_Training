package main

import (
	"fmt"
)

func main() {
	// orDone ensures that the returned channel only receives values from c until either done is closed or c is closed
	orDone := func(done, input <-chan interface{}) <-chan interface{} {
		output := make(chan interface{})
		go func() {
			defer close(output)
			for {
				select {
				case <-done:
					return
				case v, ok := <-input:
					if ok == false {
						return
					}
					select {
					case output <- v:
					case <-done:
					}
				}
			}
		}()
		return output
	}

	// bridge orchestrates multiple input channels and sends their values to a single output channel
	bridge := func(
		done <-chan interface{},
		inputStreams ...<-chan interface{},
	) <-chan interface{} {
		output := make(chan interface{})
		go func() {
			defer close(output)
			for _, inputStream := range inputStreams {
				for value := range orDone(done, inputStream) {
					select {
					case output <- value:
					case <-done:
						return
					}
				}
			}
		}()
		return output
	}

	// Use bridge to receive values from multiple streams and print them
	for value := range bridge(nil, genStreams(), genStreamsC(), genStreams()) {
		fmt.Printf("%v ", value)
	}
}

// genStreams generates a sequence of channels, each containing a single integer value
func genStreams() <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for i := 0; i < 5; i++ {
			stream <- i
		}
	}()
	return stream
}

func genStreamsC() <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for i := 10; i < 15; i++ {
			stream <- i
		}
	}()
	return stream
}
