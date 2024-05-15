package main

import (
    "errors"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// CircuitBreaker represents a circuit breaker pattern
type CircuitBreaker struct {
    mutex            sync.Mutex
    consecutiveFails int
    maxFailures      int
    cooldownDuration time.Duration
    tripped          bool
}

// NewCircuitBreaker creates a new CircuitBreaker instance
func NewCircuitBreaker(maxFailures int, cooldownDuration time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:      maxFailures,
        cooldownDuration: cooldownDuration,
    }
}

// Execute wraps the provided function with circuit breaker logic
func (cb *CircuitBreaker) Execute(req func() error) error {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if cb.tripped {
        return errors.New("circuit breaker tripped")
    }

    if cb.consecutiveFails >= cb.maxFailures {
        cb.tripped = true
        go cb.resetAfterCooldown()
        return errors.New("circuit breaker tripped")
    }

    err := req()
    if err != nil {
        cb.consecutiveFails++
    } else {
        cb.consecutiveFails = 0
    }

    return err
}

func (cb *CircuitBreaker) resetAfterCooldown() {
    time.Sleep(cb.cooldownDuration)
    cb.mutex.Lock()
    cb.tripped = false
    cb.consecutiveFails = 0
    cb.mutex.Unlock()
}

func main() {
    // Example usage of circuit breaker
    cb := NewCircuitBreaker(3, 5*time.Second) // Allow 3 consecutive failures before tripping, cooldown for 5 seconds

    // Example HTTP request function
    httpRequest := func() error {
        resp, err := http.Get("https://cbseacademic.nic.in/web_material/Notifications/2024/44_Notification_2024.pdf")
        if err != nil {
            return err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return errors.New("unexpected status code")
        }

        return nil
    }

    // Perform HTTP requests with circuit breaker protection
    for i := 0; i < 10; i++ {
        fmt.Printf("Attempt %d: ", i+1)
        err := cb.Execute(httpRequest)
        if err != nil {
            fmt.Println("Failed:", err)
        } else {
            fmt.Println("Success")
        }
        time.Sleep(1 * time.Second) // Simulate delay between requests
    }
}
