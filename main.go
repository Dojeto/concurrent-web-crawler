package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Crawler struct {
	visitedURLs map[string]bool
	mu          sync.Mutex
	wg          sync.WaitGroup
}

func (c *Crawler) fetchURL(ctx context.Context, url string, results chan<- string) {
	defer c.wg.Done()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to fetch URL:", err)
		return
	}
	defer resp.Body.Close()

	results <- fmt.Sprintf("Fetched %s with status %s", url, resp.Status)
}

func (c *Crawler) Crawl(urls []string) {
	results := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, url := range urls {
		c.wg.Add(1)
		go c.fetchURL(ctx, url, results)
	}

	go func() {
		c.wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

func main() {
	crawler := Crawler{visitedURLs: make(map[string]bool)}
	urls := []string{"https://example.com", "https://golang.org", "https://google.com"}
	start := time.Now()
	crawler.Crawl(urls)
	end := time.Since(start)
	fmt.Println(end)
}
