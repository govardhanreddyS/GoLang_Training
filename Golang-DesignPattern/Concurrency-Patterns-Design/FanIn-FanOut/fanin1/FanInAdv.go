package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Data represents the structure of data to be sent
type Data struct {
	ID   int
	Info string
}

// Service represents a service that processes data
type Service struct {
	ID int
}

// Process simulates processing of data by a service
func (s *Service) Process(ctx context.Context, data Data) <-chan Data {
	out := make(chan Data)
	go func() {
		defer close(out)
		select {
		case <-ctx.Done():
			return // Return if context is cancelled
		default:
			// Simulate processing
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			// Send back processed data
			out <- Data{ID: data.ID, Info: fmt.Sprintf("Processed by Service %d", s.ID)}
		}
	}()
	return out
}

// FanInServiceBroker combines multiple service channels into a single channel
func FanInServiceBroker(ctx context.Context, services ...*Service) <-chan Data {
	out := make(chan Data)
	for _, service := range services {
		go func(s *Service) {
			for {
				select {
				case <-ctx.Done():
					return // Return if context is cancelled
				default:
					// Receive data from service channel
					data := <-s.Process(ctx, Data{ID: rand.Intn(100), Info: "Sample Data"})
					// Send processed data to the output channel
					out <- data
				}
			}
		}(service)
	}
	return out
}

func main() {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create multiple services
	service1 := &Service{ID: 1}
	service2 := &Service{ID: 2}
	service3 := &Service{ID: 3}

	// Combine service channels into a single channel
	serviceCh := FanInServiceBroker(ctx, service1, service2, service3)

	// Read values from the combined channel
	for i := 0; i < 10; i++ {
		data := <-serviceCh
		fmt.Printf("Received: %+v\n", data)
	}

	// Cancel the context to stop processing
	cancel()
}
