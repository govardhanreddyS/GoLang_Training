package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// fetchDataFromFile simulates fetching data from a local file
func fetchDataFromFile() string {
	// Simulate reading data from a file
	return "Data from file"
}

// fetchDataFromDatabase simulates fetching data from a database
func fetchDataFromDatabase() string {
	// Simulate querying data from a database
	return "Data from database"
}

// fetchDataFromHTTP simulates fetching data from an HTTP endpoint
func fetchDataFromHTTP() string {
	// Simulate fetching data from an HTTP endpoint
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("Error fetching data from HTTP:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response:", err)
		return ""
	}

	return string(body)
}

func main() {
	// Create a wait group to wait for all worker goroutines to finish
	var wg sync.WaitGroup

	// Create channels to receive data from each source
	fileDataChan := make(chan string)
	dbDataChan := make(chan string)
	httpDataChan := make(chan string)

	// Start worker goroutines for each data source
	wg.Add(3)
	go func() {
		defer wg.Done()
		fileDataChan <- fetchDataFromFile()
	}()
	go func() {
		defer wg.Done()
		dbDataChan <- fetchDataFromDatabase()
	}()
	go func() {
		defer wg.Done()
		httpDataChan <- fetchDataFromHTTP()
	}()

	// Close channels after all data has been sent
	go func() {
		wg.Wait()
		close(fileDataChan)
		close(dbDataChan)
		close(httpDataChan)
	}()

	// Aggregate and process data received from channels
	for {
		select {
		case data, ok := <-fileDataChan:
			if !ok {
				fileDataChan = nil
			} else {
				fmt.Println("Data from file:", data)
			}
		case data, ok := <-dbDataChan:
			if !ok {
				dbDataChan = nil
			} else {
				fmt.Println("Data from database:", data)
			}
		case data, ok := <-httpDataChan:
			if !ok {
				httpDataChan = nil
			} else {
				fmt.Println("Data from HTTP:", data)
			}
		}

		// Check if all channels are closed
		if fileDataChan == nil && dbDataChan == nil && httpDataChan == nil {
			break
		}
	}

	fmt.Println("All data processed.")
}
