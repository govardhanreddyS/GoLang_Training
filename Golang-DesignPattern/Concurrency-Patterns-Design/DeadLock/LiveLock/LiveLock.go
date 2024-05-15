package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	data := make(map[string]int)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			data["key"] = data["key"] + 1
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			data["key"] = data["key"] + 1
		}
	}()

	wg.Wait()
	fmt.Println("Final value of data[\"key\"]:", data["key"])
}
