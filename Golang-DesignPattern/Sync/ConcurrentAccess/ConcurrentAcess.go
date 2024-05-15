package main

import (
	"fmt"
)

func main() {
	m := map[int]int{}
	go writeLoop(m)
	go readLoop(m)

	// stop program from exiting, must be killed
	block := make(chan struct{})
	<-block
}

func writeLoop(m map[int]int) {
	for {
		for i := 0; i < 100; i++ {
			m[i] = i
		}
	}
}

func readLoop(m map[int]int) {
	for {
		for k, v := range m {
			fmt.Println(k, "-", v)
		}
	}
}