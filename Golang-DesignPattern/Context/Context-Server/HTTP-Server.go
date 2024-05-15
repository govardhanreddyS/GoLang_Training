package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func handler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(2*time.Second))
    defer cancel()

    select {
    case <-time.After(3 * time.Second):
        fmt.Fprintln(w, "Task completed successfully")
    case <-ctx.Done():
        fmt.Fprintln(w, "Task cancelled or deadline exceeded")
    }
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
