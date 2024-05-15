package main

import (
    "context"
    "fmt"
)

type key int

func main() {
    // `ctx := context.WithValue(context.Background(), key(1), "value")` is creating a new context
	// `ctx` derived from the background context with an additional key-value pair. The key `key(1)` is
	// associated with the value `"value"` in this context. This allows you to store and retrieve
	// values associated with specific keys in the context.
	// `ctx := context.WithValue(context.Background(), key(1), "value")` is creating a new context
	// `ctx` derived from the background context with an additional key-value pair. The key `key(1)` is
	// associated with the value `"value"` in this new context. This allows you to store and retrieve
	// values associated with specific keys in the context.
	ctx := context.WithValue(context.Background(), key(1), "value")

    if v := ctx.Value(key(1)); v != nil {
        fmt.Println("Value found in context:", v)
    } else {
        fmt.Println("Value not found in context")
    }
}
