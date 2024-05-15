package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	urls := []string{
		"https://cbseacademic.nic.in/web_material/Circulars/2024/38_Circular_2024.pdf",
		"https://cbseacademic.nic.in/web_material/Circulars/2024/37_Circular_2024.pdf",
		"https://cbseacademic.nic.in/web_material/Circulars/2024/34_Circular_2024.pdf",
	}

	// Number of workers
	numWorkers := 3

	// Channel to receive tasks
	taskCh := make(chan string, len(urls))

	// Channel to indicate completion of tasks
	doneCh := make(chan bool)

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range taskCh {
				err := downloadFile(url)
				if err != nil {
					fmt.Printf("Error downloading %s: %v\n", url, err)
				} else {
					fmt.Printf("Downloaded %s successfully\n", url)
				}
			}
		}()
	}

	// Send tasks to workers
	go func() {
		for _, url := range urls {
			taskCh <- url
		}
		close(taskCh)
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	// Wait for completion
	<-doneCh
	fmt.Println("All downloads completed.")
}

func downloadFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	filename := filepath.Join("D:\\download", filepath.Base(url))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
