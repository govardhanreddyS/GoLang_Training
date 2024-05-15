package main

import (
	"fmt"
	"sync"
)

func workerA(taskA chan <-int, taskB  <- chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Perform task A
	fmt.Println("Worker A performing task A...")
	taskA <- 1
	fmt.Println("Worker A completed task A.")

	// Wait for task B to be completed
	fmt.Println("Worker A waiting for task B to complete...")
	<-taskB
	fmt.Println("Worker A completed task B.")
}

func workerB(taskA <-chan int, taskB chan <-int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Perform task B
	<-taskA
	fmt.Println("Worker B performing task B...")
	
	fmt.Println("Worker B completed task B.")

	// Wait for task A to be completed
	taskB<-1
	fmt.Println("Worker B waiting for task A to complete...")

	fmt.Println("Worker B completed task A.")
}

func main() {
	taskA := make(chan int)
	taskB := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	// Worker A
	go workerA(taskA, taskB, &wg)

	// Worker B
	go workerB(taskA, taskB, &wg)

	wg.Wait()
	fmt.Println("Both workers finished successfully.")
}
