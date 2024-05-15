package main

import (
    "testing"
)

const (
    size = 10000 // Change the size as needed
)

func BenchmarkArrayInsertion(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var arr []int
        for j := 0; j < size; j++ {
            arr = append(arr, j)
        }
    }
}

func BenchmarkMapInsertion(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m := make(map[int]bool)
        for j := 0; j < size; j++ {
            m[j] = true
        }
    }
}

func BenchmarkArrayDeletion(b *testing.B) {
    var arr []int
    for j := 0; j < size; j++ {
        arr = append(arr, j)
    }
    b.ResetTimer()
    
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(arr); j++ {
			if j == size-1 {
				arr = arr[:size-1] // Remove the last element
			} else {
				arr = append(arr[:j], arr[j+1:]...) // Remove element at index j
			}
		}
	}
	
	

}

func BenchmarkMapDeletion(b *testing.B) {
    m := make(map[int]bool)
    for j := 0; j < size; j++ {
        m[j] = true
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for j := 0; j < size; j++ {
            delete(m, j)
        }
    }
}

//go test -benchmem  -bench .