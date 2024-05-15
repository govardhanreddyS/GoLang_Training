//  Deterministically quit goroutine with quit channel option in select

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	quit := make(chan string)
	ch := generator("Hi!", quit)
	for i := rand.Intn(50); i >= 0; i-- {
		fmt.Println(<-ch, i)
	}
	quit <- "Bye!"
	fmt.Printf("Generator says %s", <-quit)
}

func generator(msg string, quit chan string) <-chan string { // returns receive-only channel
	ch := make(chan string)
	go func() { // anonymous goroutine
		run:=true
		for run {
			select {
			case ch <- fmt.Sprintf("%s", msg):
				// nothing
			case <-quit: 
			    run=false
				quit <- "See you!"
				return
			}
		}
	}()
	return ch
}
