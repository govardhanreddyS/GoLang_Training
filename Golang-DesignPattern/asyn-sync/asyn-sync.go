package main

import (
	"fmt"
	"sync"
	"time"
)

// Synchronous component
type Synchronous struct {
	mu       sync.Mutex
	data     int
	callback func(int)
}

// Asynchronous component
type Asynchronous struct {
	ch chan int
}

// NewSynchronous creates a new synchronous component
func NewSynchronous(callback func(int)) *Synchronous {
	return &Synchronous{
		callback: callback,
	}
}

// NewAsynchronous creates a new asynchronous component
func NewAsynchronous() *Asynchronous {
	return &Asynchronous{
		ch: make(chan int),
	}
}

// SetData sets the data synchronously
func (s *Synchronous) SetData(data int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = data
}

// ProcessData asynchronously processes the data
func (a *Asynchronous) ProcessData(s *Synchronous) {
	go func() {
		// Retrieve data asynchronously
		s.mu.Lock()
		data := s.data
		s.mu.Unlock()

		// Simulate processing time
		time.Sleep(time.Second)

		// Send processed data through channel
		a.ch <- data * 2
	}()
}

// GetData returns the processed data asynchronously
func (a *Asynchronous) GetData() int {
	// Retrieve processed data from channel asynchronously
	return <-a.ch
}

func main() {
	// Create synchronous and asynchronous components
	synchronous := NewSynchronous(func(data int) {
		fmt.Println("Processed data:", data)
	})
	asynchronous := NewAsynchronous()

	// Set data synchronously
	synchronous.SetData(10)

	// Process data asynchronously
	asynchronous.ProcessData(synchronous)

	// Retrieve processed data asynchronously
	processedData := asynchronous.GetData()

	// Callback function for processed data
	synchronous.callback(processedData)

	// Wait for processing to finish
	time.Sleep(2 * time.Second)
}
