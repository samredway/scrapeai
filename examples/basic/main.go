package main

import (
	"fmt"
	"log"

	"github.com/samredway/scrapeai/scrapeai"
	"github.com/samredway/scrapeai/scraping"
)

func main() {
	url := "https://example.com"

	// Example 1: Static HTML scraping
	fmt.Println("=== Static HTML Example ===")
	req, err := scrapeai.NewScrapeAiRequest(url, "Extract the main headline",
		scrapeai.WithFetchFunc(scraping.Fetch)) // Using static HTML fetching
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	result, err := scrapeai.Scrape(req)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}
	fmt.Printf("Static HTML Result: %+v\n\n", result.Results)

	// Example 2: Dynamic HTML scraping
	fmt.Println("=== Dynamic HTML Example ===")
	dynamicReq, err := scrapeai.NewScrapeAiRequest(url, "Extract the main headline")
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	// Default chromedp fetch func is used, allowing JS rendering

	result, err = scrapeai.Scrape(dynamicReq)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}
	fmt.Printf("Dynamic HTML Result: %+v\n\n", result.Results)
}
