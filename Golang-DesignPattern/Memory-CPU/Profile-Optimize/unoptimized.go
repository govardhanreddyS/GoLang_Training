package main

import (
    "fmt"
    "os"
    "runtime/pprof"
   // "strings"
    "time"
)

func concatenateStrings() string {
    var result string
    for i := 0; i < 100000; i++ {
        result += "hello" // Concatenate "hello" 100,000 times
    }
    return result
}

func main() {
    // CPU profiling
    cpuProfile, err := os.Create("cpu_unoptimized.prof")
    if err != nil {
        fmt.Println("Error creating CPU profile:", err)
        return
    }
    defer cpuProfile.Close()
    pprof.StartCPUProfile(cpuProfile)
    defer pprof.StopCPUProfile()

    // Memory profiling
    memProfile, err := os.Create("mem_unoptimized.prof")
    if err != nil {
        fmt.Println("Error creating memory profile:", err)
        return
    }
    defer memProfile.Close()

    result := concatenateStrings()
    fmt.Println("Concatenated string length:", len(result))

    // Adding some delay to simulate workload
    time.Sleep(1 * time.Second)

    // Write memory profile
    pprof.WriteHeapProfile(memProfile)
}
