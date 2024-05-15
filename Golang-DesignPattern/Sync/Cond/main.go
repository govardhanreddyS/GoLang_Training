package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	id int
}

type TaskQueue struct {
	tasks    []Task
	front    int
	rear     int
	len      int
	capacity int
}

func NewTaskQueue(capacity int) *TaskQueue {
	return &TaskQueue{
		tasks:    make([]Task, capacity),
		front:    0,
		rear:     -1,
		len:      0,
		capacity: capacity,
	}
}

func (q *TaskQueue) Enqueue(task Task) bool {
	if q.len == q.capacity {
		return false
	}
	q.rear = (q.rear + 1) % q.capacity
	q.tasks[q.rear] = task
	q.len++
	return true
}

func (q *TaskQueue) Dequeue() (Task, bool) {
	if q.len == 0 {
		return Task{}, false
	}
	task := q.tasks[q.front]
	q.front = (q.front + 1) % q.capacity
	q.len--
	return task, true
}

func main() {
	lock := sync.Mutex{}
	fullCond := sync.NewCond(&lock)
	emptyCond := sync.NewCond(&lock)
	queue := NewTaskQueue(10)

	producer := func() {
		for {
			task := Task{id: rand.Intn(1000)}
			lock.Lock()
			for !queue.Enqueue(task) {
				fmt.Println("Task queue is full")
				fullCond.Wait()
			}
			lock.Unlock()
			emptyCond.Signal()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		}
	}

	consumer := func(id int) {
		for {
			lock.Lock()
			var task Task
			for {
				var ok bool
				if task, ok = queue.Dequeue(); !ok {
					fmt.Println("Task queue is empty")
					emptyCond.Wait()
					continue
				}
				break
			}
			lock.Unlock()
			fullCond.Signal()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			fmt.Printf("Consumer %d processed task with ID: %d\n", id, task.id)
		}
	}

	for i := 0; i < 5; i++ {
		go producer()
	}
	for i := 0; i < 7; i++ {
		go consumer(i + 1)
	}

	select {}
}
