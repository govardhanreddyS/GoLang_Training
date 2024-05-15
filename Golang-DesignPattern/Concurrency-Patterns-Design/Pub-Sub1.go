package main

import (
    "fmt"
    "sync"
    "time"
)

type Publisher struct {
    subscribers map[chan string]struct{}
    mu          sync.Mutex
}

func NewPublisher() *Publisher {
    return &Publisher{
        subscribers: make(map[chan string]struct{}),
    }
}

func (p *Publisher) Subscribe() chan string {
    p.mu.Lock()
    defer p.mu.Unlock()

    ch := make(chan string, 10) // Buffer the channel to avoid blocking publisher
    p.subscribers[ch] = struct{}{}
    return ch
}

func (p *Publisher) Unsubscribe(ch chan string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    delete(p.subscribers, ch)
    close(ch) // Close the channel to signal unsubscribe
}

func (p *Publisher) Publish(message string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    for sub := range p.subscribers {
        sub <- message
    }
}

func main() {
    publisher := NewPublisher()

    // Subscribe
    subscriber1 := publisher.Subscribe()
    subscriber2 := publisher.Subscribe()

    // Read messages from subscriber 1
    go func() {
        for {
            select {
            case msg, ok := <-subscriber1:
                if !ok {
                    fmt.Println("Subscriber 1 unsubscribed")
                    return
                }
                fmt.Println("Subscriber 1 Received:", msg)
            }
        }
    }()

    // Read messages from subscriber 2
    go func() {
        for {
            select {
            case msg, ok := <-subscriber2:
                if !ok {
                    fmt.Println("Subscriber 2 unsubscribed")
                    return
                }
                fmt.Println("Subscriber 2 Received:", msg)
            }
        }
    }()

    // Publish messages
    for i := 1; i <= 5; i++ {
        msg := fmt.Sprintf("Message %d", i)
        publisher.Publish(msg)
        fmt.Println("Published:", msg)
        time.Sleep(time.Second)
    }

    // Unsubscribe one subscriber
    publisher.Unsubscribe(subscriber1)

    // Publish more messages
    for i := 6; i <= 10; i++ {
        msg := fmt.Sprintf("Message %d", i)
        publisher.Publish(msg)
        fmt.Println("Published:", msg)
        time.Sleep(time.Second)
    }

    // Wait for subscribers to process messages
    time.Sleep(5 * time.Second)
}
