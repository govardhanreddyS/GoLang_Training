package main

import (
    "testing"
)

const (
    mapSize    = 1000000 // Change the size of the map and array as needed
    arraySize  = 1000000
)

func BenchmarkMapInsertion(b *testing.B) {
    m := make(map[int]int)
    for i := 0; i < b.N; i++ {
        for j := 0; j < mapSize; j++ {
            m[j] = j
        }
    }
}

func BenchmarkArrayInsertion(b *testing.B) {
    var arr [arraySize]int
    for i := 0; i < b.N; i++ {
        for j := 0; j < arraySize; j++ {
            arr[j] = j
        }
    }
}

func BenchmarkMapAccess(b *testing.B) {
    m := make(map[int]int)
    for i := 0; i < mapSize; i++ {
        m[i] = i
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for j := 0; j < mapSize; j++ {
            _ = m[j]
        }
    }
}

func BenchmarkArrayAccess(b *testing.B) {
    var arr [arraySize]int
    for i := 0; i < arraySize; i++ {
        arr[i] = i
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for j := 0; j < arraySize; j++ {
            _ = arr[j]
        }
    }
}
//go test -benchmem -run=^$ -bench .