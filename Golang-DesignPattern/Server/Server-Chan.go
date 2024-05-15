package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Data string `json:"data"`
}

func processData(data string, resultChan chan<- Result) {
	// Simulate some processing time
	// In a real scenario, you would perform some heavy computation or I/O operation here
	result := Result{Data: fmt.Sprintf("Processed: %s", data)}

	// Send the result through the channel
	resultChan <- result
}

func processDataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestData struct {
		Data string `json:"data"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a channel to receive the result
	resultChan := make(chan Result)

	// Start a Goroutine to process the data asynchronously
	go processData(requestData.Data, resultChan)

	// Wait for the result from the channel
	result := <-resultChan

	// Encode and send the result as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/process", processDataHandler)

	// Start the server
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	fmt.Println("Server started at http://localhost:8080")

	// Block indefinitely to keep the server running
	select {}
}

/*
Send a POST Request:
You can use PowerShell's Invoke-RestMethod cmdlet to send a POST request to the /process endpoint with some JSON data. Here's an example:
powershell

$url = "http://localhost:8080/process"
$data = @{
    "data" = "example data"
}
Invoke-RestMethod -Method Post -Uri $url -Body (ConvertTo-Json $data) -ContentType "application/json"

*/