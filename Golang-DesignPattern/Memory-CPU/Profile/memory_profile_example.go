package main
//go tool pprof .\memory_profile.prof .\memory_profile.prof
import (
    "fmt"
    "os"
    "runtime/pprof"
)

func allocateMemory() {
    var slice [][]int
    for i := 0; i < 10000; i++ {
        // Allocate 1MB of memory
        slice = append(slice, make([]int, 262144))
    }
}

func main() {
    // Create a file to write the memory profile to
    f, err := os.Create("memory_profile.prof")
    if err != nil {
        fmt.Println("Error creating memory profile:", err)
        return
    }
    defer f.Close()

    // Start memory profiling
    pprof.WriteHeapProfile(f)

    // Perform memory-intensive operation
    allocateMemory()

    // Stop memory profiling
    pprof.WriteHeapProfile(f)

    fmt.Println("Memory profiling complete. Profile data written to memory_profile.prof")
}
