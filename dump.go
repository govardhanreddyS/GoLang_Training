package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Job represents a task to be performed by a worker
type Job struct {
	ID       int
	Endpoint string
}

// Result represents the result of processing a job
type Result struct {
	JobID int
	Title string // Example: Include a field to store the title of the fetched data
}

// Worker represents a worker that performs a job
type Worker struct {
	ID         int
	jobCh      <-chan Job
	resultCh   chan<- Result
	httpClient *http.Client
}

// NewWorker creates a new worker with the given ID, job channel, result channel, and HTTP client
func NewWorker(id int, jobCh <-chan Job, resultCh chan<- Result, httpClient *http.Client) *Worker {
	return &Worker{
		ID:         id,
		jobCh:      jobCh,
		resultCh:   resultCh,
		httpClient: httpClient,
	}
}

// ProcessJob simulates processing a job by fetching data from an API endpoint and returning a result
func (w *Worker) ProcessJob(job Job) {
	// Fetch data from the API endpoint
	resp, err := w.httpClient.Get(job.Endpoint)
	if err != nil {
		fmt.Printf("Worker %d: Error fetching data from %s: %v\n", w.ID, job.Endpoint, err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Worker %d: Error reading response body from %s: %v\n", w.ID, job.Endpoint, err)
		return
	}

	// Parse the JSON response
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Worker %d: Error parsing JSON response from %s: %v\n", w.ID, job.Endpoint, err)
		return
	}

	// Extract specific fields of interest (e.g., title)
	title, ok := data["title"].(string)
	if !ok {
		fmt.Printf("Worker %d: Unable to extract title from JSON response for %s\n", w.ID, job.Endpoint)
		return
	}

	// Send the result to the result channel
	w.resultCh <- Result{JobID: job.ID, Title: title}
	fmt.Printf("Worker %d: Processed job %d, Title: %s\n", w.ID, job.ID, title)
}

func main() {
	// Define the API endpoints to fetch data from
	endpoints := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
	}

	// Create a wait group to wait for all worker goroutines to finish
	var wg sync.WaitGroup

	// Create channels for jobs and results
	jobCh := make(chan Job, len(endpoints))
	resultCh := make(chan Result, len(endpoints))

	// Create an HTTP client
	httpClient := &http.Client{}

	// Create and start worker goroutines
	numWorkers := len(endpoints)
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, jobCh, resultCh, httpClient)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				worker.ProcessJob(job)
			}
		}()
	}

	// Create and enqueue jobs
	for i, endpoint := range endpoints {
		jobCh <- Job{ID: i, Endpoint: endpoint}
	}
	close(jobCh)

	// Wait for all workers to finish processing
	wg.Wait()

	// Close the result channel
	close(resultCh)

	// Collect and print results
	for result := range resultCh {
		fmt.Printf("Result for Job %d: Title - %s\n", result.JobID, result.Title)
	}
}
