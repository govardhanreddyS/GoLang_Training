package main
//http://localhost:6060/debug/pprof/
import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // Import to enable pprof endpoints
	"sync"
	"time"
)

// hardWork simulates CPU and memory-intensive work
func hardWork(wg *sync.WaitGroup) {
	defer wg.Done() // Notify WaitGroup that this goroutine is done
	fmt.Printf("Start: %v\n", time.Now())

	// Memory-intensive operation
	a := make([]string, 0, 500000)
	for i := 0; i < 500000; i++ {
		a = append(a, "aaaa")
	}

	// Blocking operation
	time.Sleep(2 * time.Second)
	fmt.Printf("End: %v\n", time.Now())
}

func main() {
	var wg sync.WaitGroup

	// Start an HTTP server to expose pprof endpoints
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Increment WaitGroup counter for the HTTP server goroutine
	wg.Add(1)

	// Start the hardWork goroutine
	wg.Add(1)
	go hardWork(&wg)

	// Wait for both goroutines to finish
	wg.Wait()
}
