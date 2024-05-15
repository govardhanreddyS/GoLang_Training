package main

import (
    "fmt"
    "sync"
)     

func main() {
    var wg sync.WaitGroup
    pingCh := make(chan string)
    pongCh := make(chan string)   

    wg.Add(1) // One WaitGroup for the first goroutine
    go func() {
        defer wg.Done()                  
        for i := 0; i < 26; i++ {          //Ping 
            letter := string(rune('A' + i)) // Convert integer to corresponding letter
            fmt.Println("Ping", i)
            pongCh <- letter
            <-pingCh
          // println("next" , <-pingCh)
        }
    }()
    go func() { //Pong 
        wg.Add(1) // One WaitGroup for the second goroutine
        defer wg.Done()
        for i := 0; i < 26; i++ {
            letter := <-pongCh
            fmt.Println("Pong", letter)
            pingCh <- string(rune(letter[0] + 1)) // Send next letter in sequence
        }
    }()

    wg.Wait() // Wait for the first goroutine to finish
    close(pingCh)
    close(pongCh)
}