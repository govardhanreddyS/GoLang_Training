package main

import (
	"fmt"
)

// Data represents structured data with multiple fields
type Data struct {
	ID    int
	Name  string
	Value float64
}

func main() {
	// Define repeat function to generate a stream of structured data
	repeat := func(done <-chan interface{}, values ...Data) <-chan Data {
		valueStream := make(chan Data)

		go func() {
			defer close(valueStream)

			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()

		return valueStream
	}

	// Define take function to take a specified number of values from the stream
	take := func(done <-chan interface{}, valueStream <-chan Data, num int) <-chan Data {
		takeStream := make(chan Data)

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()

		return takeStream
	}

	// Define tee function to split a single input channel into two output channels
	tee := func(done <-chan interface{}, in <-chan Data) (<-chan Data, <-chan Data) {
		out1 := make(chan Data)
		out2 := make(chan Data)

		go func() {
			defer close(out1)
			defer close(out2)

			for val := range in {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()

		return out1, out2
	}

	// Create a cancellation channel
	done := make(chan interface{})
	defer close(done)

	// Create a stream of structured data
	dataStream := repeat(done, Data{ID: 1, Name: "John", Value: 10.5}, Data{ID: 2, Name: "Alice", Value: 20.7})

	// Split the stream into two channels
	out1, out2 := tee(done, take(done, dataStream, 4))

	// Read values from the output channels and print them
	for val1 := range out1 {
		val2 := <-out2
		fmt.Printf("out1: %+v, out2: %+v\n", val1, val2)
	}
}
