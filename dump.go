package main

import (
	"fmt"
	"math/rand"
	//"time"
)

// gen generates random numbers and sends them to the returned channel
func gen() <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			out <- rand.Intn(100)
		}
	}()
	return out
}

// fanIn combines multiple input channels into a single channel
func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	for _, in := range inputs {
		go func(ch <-chan int) {
			for n := range ch {
				out <- n
			}
		}(in)
	}
	return out
}

func main() {
	// Create input channels
	ch1 := gen()
	ch2 := gen()
	ch3 := gen()

	// Combine input channels into a single channel
	fanInCh := fanIn(ch1, ch2, ch3)

	// Read values from the combined channel
	for i := 0; i < 10; i++ {
		fmt.Println(<-fanInCh)
	}
}
