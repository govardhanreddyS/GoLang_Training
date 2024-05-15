package main

import (
	"fmt"
)

var repeat = func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

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

var take = func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

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

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					fmt.Println("finish")
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

	tee := func(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {

		out1 := make(chan interface{})
		out2 := make(chan interface{})

		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
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

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	for vall := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", vall, <-out2)
	}
}

/*
In This program demonstrates a tee function that splits a single input channel into two output channels, duplicating each value received from the input channel.

Here's what happens in the program:

The repeat function creates a channel valueStream that repeatedly sends the provided values. This function is intended to generate an infinite stream of values but is limited by the context of the program.

The take function creates a channel takeStream that takes a specified number of values from the valueStream channel.

The orDone function ensures that the valStream channel only receives values from the c channel until either done is closed or c is closed. It's a common pattern to prevent goroutines from blocking indefinitely when reading from closed channels.

The tee function splits the input channel in into two output channels out1 and out2. Each value received from the input channel is duplicated and sent to both output channels.

In the main function, a done channel is created for cancellation purposes.

The repeat function is used to generate a stream of values 1 and 2, and the take function is used to take 4 values from this stream.

The tee function is called with the done channel and the channel returned by take as input.

Finally, a loop reads values from out1 and out2 alternately and prints them. Since both channels receive the same values, the output will show pairs of values from the input stream.

The program makes use of defer statements to ensure proper cleanup of resources, closing channels when they are no longer needed.

When the program completes, it prints pairs of values from the two output channels until the done channel is closed, at which point it prints "finish" and terminates.

This program demonstrates concurrency patterns such as splitting channels and handling cancellation to safely manage concurrent operations.
*/