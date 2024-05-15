package main

import (
	"fmt"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func printer(in <-chan int) {
	for n := range in {
		fmt.Println(n)
	}
}

func main() {
	// Create a generator stage
	nums := []int{1, 2, 3, 4, 5}
	gen := generator(nums...)

	// Create a square stage
	sq := square(gen)

	// Print the squared numbers
	printer(sq)
}
