package main

import (
	"fmt"
	"io"
	"net/http"
)

func downloadFile(url string, ch chan<- error) {
	// Download the file at the specified URL
	resp, err := http.Get(url)
	if err != nil {
		ch <- err // Send any error encountered during download
		return
	}
	defer resp.Body.Close()

	// Simulate writing the downloaded content (replace with actual writing logic)
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		ch <- err
		return
	}

	fmt.Println("Downloaded:", url)
	ch <- nil // Send nil to indicate successful download
}

func main() {
	urls := []string{
		"https://golang.org/doc/install",
		"https://www.google.com",
		"https://github.com",
	}

	// Channel to receive errors or success signals from downloads
	errCh := make(chan error)

	// Start downloading files concurrently using goroutines
	for _, url := range urls {
		go downloadFile(url, errCh)
	}

	// Loop to handle errors or success signals from downloads
	for i := range urls {
		err := <-errCh
		if err != nil {
			fmt.Println("Error downloading", urls[i], err)
		} else {
			fmt.Println("Successfully downloaded", urls[i])
		}
	}

	fmt.Println("All downloads finished!")
}
