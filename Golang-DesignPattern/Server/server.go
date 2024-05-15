package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
)

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var (
    todos     []Todo
    todosLock sync.RWMutex
)

func getTodos(w http.ResponseWriter, r *http.Request) {
    todosLock.RLock()
    defer todosLock.RUnlock()

    json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
    var newTodo Todo
    err := json.NewDecoder(r.Body).Decode(&newTodo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todosLock.Lock()
    defer todosLock.Unlock()

    newTodo.ID = len(todos) + 1
    todos = append(todos, newTodo)

    json.NewEncoder(w).Encode(newTodo)
}

func main() {
    // Initialize some example todos
    todos = append(todos, Todo{ID: 1, Title: "Learn Go", Done: false})
    todos = append(todos, Todo{ID: 2, Title: "Build a REST API", Done: false})

    http.HandleFunc("/todos", getTodos)
    http.HandleFunc("/todos/add", addTodo)

    // Start the server in a Goroutine to handle requests concurrently
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    fmt.Println("Server started at http://localhost:8080")

    // Wait for interrupt signal to stop the server
    <-make(chan struct{})
}
