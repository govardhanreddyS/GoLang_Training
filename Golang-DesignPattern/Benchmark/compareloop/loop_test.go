package main

import (
    "math/big"
    "testing"
)

const size = 10000000 // Size of the slice

var nums []int

func init() {
    // Initialize the slice with numbers from 1 to size
    nums = make([]int, size)
    for i := 0; i < size; i++ {
        nums[i] = i + 1
    }
}

func BenchmarkLoopSum(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sum := 0
        for _, num := range nums {
            sum += num
        }
    }
}

func BenchmarkBigSum(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var sum big.Int
        for _, num := range nums {
            sum.Add(&sum, big.NewInt(int64(num)))
        }
    }
}
