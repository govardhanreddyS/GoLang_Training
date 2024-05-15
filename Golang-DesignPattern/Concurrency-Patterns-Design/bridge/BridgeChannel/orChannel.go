package main

import (
    "fmt"
    "time"
)

func orDone(done, c <-chan interface{}) <-chan interface{} {
    valStream := make(chan interface{})
    go func() {
        defer close(valStream)
        for {
            select {
            case <-done:
                return
            case v, ok := <-c:
                if ok == false {
                    return
                }
                select {
                case valStream <- v:
                case <-done:
                }
            }
        }
    }()
    return valStream
}

func main() {
    done := make(chan interface{})
    defer close(done)

    // Assume we have multiple channels receiving data
    channel1 := make(chan interface{})
    channel2 := make(chan interface{})

    // Simulate sending data to channels
    go func() {
        defer close(channel1)
        for i := 0; i < 5; i++ {
            channel1 <- i
            time.Sleep(500*time.Millisecond)
        }
    }()
    go func() {
        defer close(channel2)
        for i := 5; i < 10; i++ {
            channel2 <- i
        }
    }()

    // Aggregate data from multiple channels
    aggregatedStream := orDone(done, channel1)
    aggregatedStream1 := orDone(done, channel2)

    // Process aggregated data
    for v := range aggregatedStream {
        fmt.Printf("%v ", v)
    }

    for v := range aggregatedStream1 {
        fmt.Printf("%v ", v)
    }
}
