package main

import (
	"fmt"
	"time"
)

// ActiveObject represents the active object
type ActiveObject struct {
	// channel to receive method invocations
	requestChannel chan func()

	// channel to receive results
	resultChannel chan interface{}
}

// NewActiveObject creates a new active object
func NewActiveObject() *ActiveObject {
	obj := &ActiveObject{
		requestChannel: make(chan func(), 10), // buffered channel
		resultChannel:  make(chan interface{}),
	}

	// start the active object loop
	go obj.run()

	return obj
}

// run starts the active object loop
func (ao *ActiveObject) run() {
	for {
		// receive method invocation from the request channel
		method := <-ao.requestChannel

		// execute the method
		method()

		// send result (if any) to the result channel
		ao.resultChannel <- "Method execution completed"
	}
}

// AddTask adds a task to be executed asynchronously by the active object
func (ao *ActiveObject) AddTask(task func()) {
	// send method invocation to the request channel
	ao.requestChannel <- task
}

func main() {
	// create a new active object
	activeObj := NewActiveObject()

	// Add tasks to the active object
	for i := 1; i <= 5; i++ {
		taskID := i
		activeObj.AddTask(func() {
			fmt.Printf("Task %d is being executed\n", taskID)
			time.Sleep(1 * time.Second)
			fmt.Printf("Task %d has completed\n", taskID)
		})
	}

	// Wait for tasks to complete
	for i := 1; i <= 5; i++ {
		<-activeObj.resultChannel
	}

	fmt.Println("All tasks have completed.")
}
