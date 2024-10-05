package main

import (
	"fmt"
	"log"

	"github.com/samredway/scrapeai/scrapeai"
)

func main() {
	url := "https://example.com"

	// Example using scrapeai package to get the main headline
	req := scrapeai.ScrapeAiRequest{
		Url:    url,
		Prompt: "Extract the main headline",
	}
	result, err := scrapeai.Scrape(req)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}
	fmt.Printf("AI Scraping Result: %+v\n", result)

	// Example using scrapeai package to get the main body of the page
	req_2 := scrapeai.ScrapeAiRequest{
		Url:    url,
		Prompt: "Extract the main body of the page, exclude the main headline",
	}
	result_2, err := scrapeai.Scrape(req_2)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}
	fmt.Printf("AI Scraping Result: %+v\n", result_2)

	// NOTE that we fetch the results by default using a temparature of 0.0 and
	// a constant seed. This means that we get consistent results across
	// multiple calls.
}
