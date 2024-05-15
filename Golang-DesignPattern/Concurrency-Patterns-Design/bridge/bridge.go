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
		inputStream <-chan <-chan interface{},
	) <-chan interface{} {
		output := make(chan interface{})
		go func() {
			defer close(output)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-inputStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				for value := range orDone(done, stream) {
					select {
					case output <- value:
					case <-done:
					}
				}
			}
		}()
		return output
	}

	// genStreams generates a sequence of channels, each containing a single integer value
	genStreams := func() <-chan <-chan interface{} {
		streams := make(chan (<-chan interface{}))
		go func() {
			defer close(streams)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				streams <- stream
			}
		}()
		return streams
	}

	// Use bridge to receive values from the generated streams and print them
	for value := range bridge(nil, genStreams()) {
		fmt.Printf("%v ", value)
	}
}
