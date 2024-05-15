package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// Task represents the download and processing job
type Task struct {
	URL   string
	Data  []byte // Data field to store the downloaded content
	Error error
}

// ProcessedData represents the final processed data
type ProcessedData struct {
	URL       string
	Processed string // Replace with the actual processed data type (e.g., string, map)
}

// DownloadWorker downloads content from a URL and saves it to the specified directory
func downloadWorker(id int, tasks <-chan string, results chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range tasks {
		resp, err := http.Get(url)
		if err != nil {
			results <- Task{URL: url, Error: err}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			results <- Task{URL: url, Error: fmt.Errorf("unexpected status code: %d", resp.StatusCode)}
			resp.Body.Close()
			continue
		}

		// Create the file with the same name as in the URL
		filePath := "D://download/" + filepath.Base(url)
		file, err := os.Create(filePath)
		if err != nil {
			results <- Task{URL: url, Error: err}
			resp.Body.Close()
			continue
		}
		defer file.Close()

		// Copy the response body to the file
		_, err = io.Copy(file, resp.Body)
		resp.Body.Close()
		if err != nil {
			results <- Task{URL: url, Error: err}
			continue
		}

		fmt.Printf("Worker %d downloaded content from %s and saved it to %s\n", id, url, filePath)

		// Read the response body and store it in the Data field of Task
		data, err := io.ReadAll(file)
		if err != nil {
			results <- Task{URL: url, Error: err}
			continue
		}

		// Store the downloaded content in the Task struct and send it over the results channel
		results <- Task{URL: url, Data: data}
	}
}

// ProcessWorker processes downloaded data and stores the result in a slice
func processWorker(id int, tasks <-chan Task, processedData *[]ProcessedData, wg *sync.WaitGroup) {
	defer wg.Done() // Call wg.Done() after processing

	for task := range tasks {
		if task.Error != nil {
			fmt.Printf("Error downloading %s: %v\n", task.URL, task.Error)
			continue
		}

		// Process data (replace with your actual processing logic)
		processed := string(task.Data) // Convert byte slice to string

		// Append processed data to the slice
		*processedData = append(*processedData, ProcessedData{URL: task.URL, Processed: processed})

		fmt.Printf("Worker %d processed data from %s\n", id, task.URL)
	}
}

func main() {
	var wg sync.WaitGroup

	// Define number of workers
	numDownloadWorkers := 10
	numProcessWorkers := 3

	// Create channels for tasks and results
	downloadTasks := make(chan string)
	downloadResults := make(chan Task)
	processedData := make([]ProcessedData, 0) // Slice to store processed data

	// Launch download worker pool
	for i := 0; i < numDownloadWorkers; i++ {
		wg.Add(1)
		go downloadWorker(i, downloadTasks, downloadResults, &wg)
	}

	// Launch process worker pool
	processWG := sync.WaitGroup{} // New wait group for processing workers
	processWG.Add(numProcessWorkers)
	for i := 0; i < numProcessWorkers; i++ {
		go processWorker(i, downloadResults, &processedData, &processWG) // Corrected code
	}

	// Define URLs to download (replace with your actual URLs)
	urls := []string{
		"https://cbseacademic.nic.in/web_material/Circulars/2024/38_Circular_2024.pdf",
		"https://cbseacademic.nic.in/web_material/Circulars/2024/37_Circular_2024.pdf",
		"https://cbseacademic.nic.in/web_material/Circulars/2024/34_Circular_2024.pdf",
	}

	// Send download tasks
	for _, url := range urls {
		downloadTasks <- url
	}

	// Close the download tasks channel to signal workers no more tasks are coming
	close(downloadTasks)

	// Wait for all download workers to finish
	go func() {
		wg.Wait()
		// Close the download results channel after all download workers are done
		close(downloadResults)
	}()

	// Wait for all processing workers to finish
	processWG.Wait() // Corrected code

	// Print processed data (if needed)
	fmt.Println("Processed data:")
	for _, data := range processedData {
		fmt.Printf("  - URL: %s, Processed: %s\n", data.URL, data.Processed)
	}

	fmt.Println("All tasks processed!")
}
