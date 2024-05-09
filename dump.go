package main

import (
    //"math"
    "testing"
)

func BenchmarkFactorialIterative(b *testing.B) {
    for i := 0; i < b.N; i++ {
        factorialIterative(10) // Change the input value as needed
    }
} //go test -benchmem -run=^$ -bench ^BenchmarkFactorialIterative$ bench-test
//go test -benchmem -bench .  

func BenchmarkFactorialRecursive(b *testing.B) {
    for i := 0; i < b.N; i++ {
        factorialRecursive(10) // Change the input value as needed
    }
}

// Function to calculate factorial iteratively
func factorialIterative(n int) int {
    result := 1
    for i := 1; i <= n; i++ {
        result *= i
    }
    return result
}

// Function to calculate factorial recursively
func factorialRecursive(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorialRecursive(n-1)
}
//go test -benchmem -bench .          
/*
BenchmarkFactorialIterative-12: This benchmark ran 167881663 iterations of the
 factorialIterative function. It took an average of 7.221 nanoseconds per operation, 
 and there were 0 memory allocations (B/op) and 0 allocations per operation (allocs/op).
BenchmarkFactorialRecursive-12: This benchmark ran 81732176 iterations of the factorialRecursive
 function. It took an average of 14.43 nanoseconds per operation, and there were 0 memory 
 allocations (B/op) and 0 allocations per operation (allocs/op).
Both benchmarks passed (PASS), indicating that the functions performed as expected.
This output provides insights into the performance of both iterative and recursive approaches to
 calculating factorials. In this case, the iterative approach (factorialIterative) appears to be faster 
 than the recursive approach (factorialRecursive).
*/
