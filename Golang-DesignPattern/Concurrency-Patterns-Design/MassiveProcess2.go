/*
Concurrent Map-Reduce with Structured Data:
Let's say you have a large dataset represented by structured data, 
and you want to perform a map-reduce operation on it concurrently.
*/

package main

import (
    "fmt"
    "sync"
)

// Data represents structured data for map-reduce
type Data struct {
    Key   string
    Value int
}

func main() {
    data := []Data{
        {"A", 1},
        {"B", 2},
        {"C", 3},
        {"A", 4},
        {"B", 5},
    }

    // Map function
    mapper := func(d Data) (string, int) {
        return d.Key, d.Value * 2
    }

    // Reduce function
    reducer := func(key string, values []int) int {
        sum := 0
        for _, v := range values {
            sum += v
        }
        return sum
    }

    // Concurrent map-reduce
    var (
        result   = make(map[string]int)
        resultMu sync.Mutex
        wg       sync.WaitGroup
    )

    // Map phase
    for _, d := range data {
        wg.Add(1)
        go func(d Data) {
            defer wg.Done()
            key, value := mapper(d)
            resultMu.Lock()
            result[key] += value
            resultMu.Unlock()
        }(d)
    }
    wg.Wait()

    // Reduce phase
    for key, values := range result {
        fmt.Printf("%s: %d\n", key, reducer(key, []int{values}))
    }
}
