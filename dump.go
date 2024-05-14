/*
The orDone function in the provided code is a utility function that ensures graceful termination of a channel when a done channel is closed. It's particularly useful in scenarios where you have multiple goroutines reading from the same channel, and you need to signal all goroutines to stop processing when a termination condition is met.

Here's a brief explanation of how orDone works and a potential use case:

How orDone Works:
orDone takes two input channels: done and c.
It returns a new channel, valStream, where values from c are sent.
Inside a goroutine, it continuously reads from c and forwards the values to valStream.
It also listens to the done channel. If done is closed, it stops processing and closes valStream.
Use Case:
Suppose you have a scenario where you're receiving data from multiple sources through channels, and you want to aggregate this data into a single stream. Additionally, you want to gracefully handle the termination of the operation when a termination signal is received. Here's how you could use orDone in such a scenario:


*/
package main

import (
	"fmt"
)

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func genVals() <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	for v := range bridge(done, genVals()) {
		fmt.Printf("%v ", v)
	}
}
