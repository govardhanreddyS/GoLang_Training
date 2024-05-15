package main

import (
    "math/rand"
    "sort"
    "testing"
    "time"
)

const sizex = 1000 // Size of the slice

var nums []int

func init() {
    // Initialize the slice with random numbers
    rand.Seed(time.Now().UnixNano())
    nums = rand.Perm(sizex)
}

func bubbleSort(nums []int) {
    n := len(nums)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if nums[j] > nums[j+1] {
                nums[j], nums[j+1] = nums[j+1], nums[j]
            }
        }
    }
}

func BenchmarkSort(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sorted := make([]int, len(nums))
        copy(sorted, nums)
        sort.Ints(sort.IntSlice(sorted))
    }
}

func BenchmarkBubbleSort(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sorted := make([]int, len(nums))
        copy(sorted, nums)
        bubbleSort(sorted)
    }
}
