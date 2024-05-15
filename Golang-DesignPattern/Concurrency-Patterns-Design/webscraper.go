package main

import (
	"fmt"
	"net/http"
	"sync"
)

// URLList represents a list of URLs to scrape
var URLList = []string{
	"https://example.com/page1",
	"https://example.com/page2",
	"https://example.com/page3",
	"https://example.com/page4",
	"https://example.com/page5",
}

// Result represents the result of a scraping operation
type Result struct {
	URL     string
	Content string
}

// Scraper scrapes the content of a URL and sends the result to the output channel
func Scraper(url string, out chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %s\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	content := ""
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			content += string(buf[:n])
		}
		if err != nil {
			break
		}
	}

	// Send the result to the output channel
	out <- Result{URL: url, Content: content}
}

func main() {
	// Create channels
	results := make(chan Result)

	// Create a wait group
	var wg sync.WaitGroup

	// Start scraper workers
	for _, url := range URLList {
		wg.Add(1)
		go Scraper(url, results, &wg)
	}

	// Start a goroutine to close the results channel once all scrapers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and process results
	for result := range results {
		fmt.Printf("Scraped %s: %d bytes\n", result.URL, len(result.Content))
	}
}
