package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	data = make(map[string]int)
	mu   sync.RWMutex
)

func main() {
	// Start writer goroutine to periodically update the map
	go writer()

	// Start multiple reader goroutines
	for i := 0; i < 3; i++ {
		go reader(i)
	}

	// Wait indefinitely
	select {}
}

func writer() {
	for {
		// Acquire write lock
		mu.Lock()

		// Update map
		data["value"]++
		fmt.Println("Writer updated map with value:", data["value"])

		// Release write lock
		mu.Unlock()

		// Sleep for some time
		time.Sleep(100* time.Microsecond)
	}
}

func reader(id int) {
	for {
		// Acquire read lock
		mu.RLock()

		// Read value from map
		value := data["value"]
		fmt.Printf("Reader %d read value: %d\n", id, value)

		// Release read lock
		mu.RUnlock()

		// Sleep for some time
		time.Sleep(100* time.Microsecond)
	}
}
