package main

import (
    "fmt"
    "sync"
    "time"
)

type Publisher struct {
    subscribers map[chan Data]struct{}
    mu          sync.Mutex
}

type Data struct {
    Message string
    // Add more fields as needed
}

func NewPublisher() *Publisher {
    return &Publisher{
        subscribers: make(map[chan Data]struct{}),
    }
}

func (p *Publisher) Subscribe() chan Data {
    p.mu.Lock()
    defer p.mu.Unlock()

    ch := make(chan Data, 10) // Buffer the channel to avoid blocking publisher
    p.subscribers[ch] = struct{}{}
    return ch
}

func (p *Publisher) Unsubscribe(ch chan Data) {
    p.mu.Lock()
    defer p.mu.Unlock()

    delete(p.subscribers, ch)
    close(ch) // Close the channel to signal unsubscribe
}

func (p *Publisher) Publish(data Data) {
    p.mu.Lock()
    defer p.mu.Unlock()

    for sub := range p.subscribers {
        sub <- data
    }
}

func main() {
    publisher := NewPublisher()

    numSubscribers := 3
    subscribers := make([]chan Data, numSubscribers)

    // Subscribe
    for i := 0; i < numSubscribers; i++ {
        subscribers[i] = publisher.Subscribe()
    }

    // Read messages from subscribers
    for i, sub := range subscribers {
        go func(idx int, ch chan Data) {
            for {
                select {
                case data, ok := <-ch:
                    if !ok {
                        fmt.Printf("Subscriber %d unsubscribed\n", idx)
                        return
                    }
                    fmt.Printf("Subscriber %d Received: %+v\n", idx, data)
                }
            }
        }(i+1, sub)
    }

    // Publish messages
    for i := 1; i <= 5; i++ {
        data := Data{
            Message: fmt.Sprintf("Message %d", i),
            // Add more fields as needed
        }
        publisher.Publish(data)
        fmt.Println("Published:", data)
        time.Sleep(time.Second)
    }

    // Unsubscribe one subscriber
    publisher.Unsubscribe(subscribers[0])

    // Publish more messages
    for i := 6; i <= 10; i++ {
        data := Data{
            Message: fmt.Sprintf("Message %d", i),
            // Add more fields as needed
        }
        publisher.Publish(data)
        fmt.Println("Published:", data)
        time.Sleep(time.Second)
    }

    // Wait for subscribers to process messages
    time.Sleep(5 * time.Second)
}
