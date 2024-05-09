package main

import (
    "fmt"
    "sync"
    "time"
)

func routine(id int, wg *sync.WaitGroup, ch chan string, pauseCh <-chan bool) {
    defer wg.Done()
    for {
        select {
        case <-pauseCh:
            fmt.Printf("Routine %d paused\n", id)
            <-pauseCh // Wait for unpause signal
            fmt.Printf("Routine %d resumed\n", id)
        default:
            fmt.Printf("Routine %d running\n", id)
            // Simulate some work
            time.Sleep(time.Second*20)
            ch <- fmt.Sprintf("Routine %d completed", id)
        }
    }
}

func main() {
    numRoutines := 3
    ch := make(chan string, numRoutines)
    pauseCh := make(chan bool)

    var wg sync.WaitGroup
    wg.Add(numRoutines)

    // Start routines
    for i := 1; i <= numRoutines; i++ {
        go routine(i, &wg, ch, pauseCh)
    }

    // Control loop
    for {
        // Round-robin scheduling
        select {
        case result := <-ch:
            fmt.Println(result)
        case <-time.After(5 * time.Second):
            // Pause all routines for 2 seconds
            fmt.Println("Pausing all routines...")
            for i := 0; i < numRoutines; i++ {
                pauseCh <- true
            }
            time.Sleep(3 * time.Second)
            // Resume all routines
            fmt.Println("Resuming all routines...")
            for i := 0; i < numRoutines; i++ {
                pauseCh <- false
            }
        }
    }

    // Wait for routines to finish (which will never happen in this example)
    wg.Wait()
}
