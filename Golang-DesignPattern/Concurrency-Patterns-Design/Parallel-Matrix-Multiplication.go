package main

import (
    "fmt"
    "sync"
)

func multiply(a, b [][]int, result chan<- [][]int, wg *sync.WaitGroup) {
    defer wg.Done()
    rowsA, colsA := len(a), len(a[0])
    colsB := len(b[0])
    res := make([][]int, rowsA)
    for i := range res {
        res[i] = make([]int, colsB)
    }
    for i := 0; i < rowsA; i++ {
        for j := 0; j < colsB; j++ {
            for k := 0; k < colsA; k++ {
                res[i][j] += a[i][k] * b[k][j]
            }
        }
    }
    result <- res
}

func main() {
    a := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
    b := [][]int{{9, 8, 7}, {6, 5, 4}, {3, 2, 1}}
    result := make(chan [][]int)
    var wg sync.WaitGroup

    for i := 0; i < len(a); i++ {
        for j := 0; j < len(b[0]); j++ {
            wg.Add(1)
            go multiply(a, b, result, &wg)
        }
    }

    go func() {
        wg.Wait()
        close(result)
    }()

    res := <-result
    fmt.Println("Result:", res)
}
