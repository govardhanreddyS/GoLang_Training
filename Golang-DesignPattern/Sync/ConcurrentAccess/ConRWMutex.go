package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	data map[string]string
	lock sync.RWMutex
)

func main() {
	data = make(map[string]string)
	wticker := time.NewTicker(time.Millisecond * 500)
	rticker := time.NewTicker(time.Millisecond * 500)

	
	//Two goroutines can read at the same time
	go func() {
		for t := range rticker.C {
			lock.RLock()
			time.Sleep(time.Millisecond * 1500)
			fmt.Println("read 1 at", t, len(data))
			lock.RUnlock()
		}
	}()
	go func() {
		for t := range rticker.C {
			lock.RLock()
			time.Sleep(time.Millisecond * 1500)
			fmt.Println("read 2 at", t, len(data))
			lock.RUnlock()
		}
	}()
	time.Sleep(5 * time.Second)

	
	// After this goroutine starts, the read activity will be stuck.
	go func() {
		counter := 0
		
		// There is a read operation between the following two time points. Write operations cannot be performed during the read operation.
		fmt.Println("request write data at ", time.Now())
		lock.Lock()
		fmt.Println("start write data at ", time.Now())
		for t := range wticker.C {
			counter++
			fmt.Println("write at", t)
			data[strconv.Itoa(counter)] = strconv.Itoa(counter)
			if counter > 10 {
				break
			}
		}
		lock.Unlock()
	}()
	select {}
}