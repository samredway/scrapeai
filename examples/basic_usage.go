package main

import (
	"fmt"
	"log"

	"github.com/samredway/scrapeai/scrapeai"
	"github.com/samredway/scrapeai/scraping"
)

func main() {
	url := "https://example.com"

	// Example using scrapeai package to get the main headline
	// Here we use scraping.Fetch as a the fetch func which is a naive fetch
	// that just gets the static html from the url, this is a little faster
	// than using the default chromedp fetch func
	req := scrapeai.NewScrapeAiRequest(url, "Extract the main headline", scrapeai.WithFetchFunc(scraping.Fetch))
	result, err := scrapeai.Scrape(req)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}
	fmt.Printf("AI Scraping Result: %+v\n", result)

	// Example using scrapeai package to get the main body of the page
	// Here we use the default chromedp fetch func which allows the page to be
	// fetched dynamically with all the javascript and css rendered
	req2 := scrapeai.NewScrapeAiRequest(url, "Extract the main headline")
	result2, err := scrapeai.Scrape(req2)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}
	fmt.Printf("AI Scraping Result: %+v\n", result2)

	// NOTE that we fetch the results by default using a temparature of 0.0 and
	// a constant seed. This means that we get consistent results across
	// multiple calls.
}
