package main

import (
	"context"
	"fmt"
	"log"

	"github.com/samredway/scrapeai/scrapeai"
	"github.com/samredway/scrapeai/scraping"
)

// Why use ScrapeAI with Zyte instead of Zyte's own AI scrapers?
//
// 1. Cost efficiency: ScrapeAI uses OpenAI's API directly, which is typically
//    more cost-effective than Zyte's AI scraping services that include markup.
//    You pay OpenAI's rates plus Zyte's proxy costs, rather than Zyte's bundled
//    AI scraping pricing.
//
// 2. Granular control: With ScrapeAI, you have full control over:
//    - Custom prompts tailored to your specific extraction needs
//    - Custom JSON schemas for structured data output
//    - Model selection (gpt-4o-mini, gpt-4, etc.)
//    - Temperature, seed, and other model parameters
//
// 3. Flexibility: Easily switch between different fetch methods (static HTML,
//    chromedp for JS rendering, or Zyte proxy for anti-bot protection) while
//    maintaining the same AI extraction interface.
//
// 4. Integration: Built as a Go library, ScrapeAI integrates seamlessly into
//    your existing Go applications without requiring external API calls to
//    Zyte's AI scraping endpoints.

func main() {
	ctx := context.Background()
	url := "https://example.com"

	// Example: Using Zyte Static Proxy
	fmt.Println("=== Zyte Static Proxy Example ===")
	fmt.Println()

	req, err := scrapeai.NewScrapeAiRequest(url, "Extract the main headline",
		scrapeai.WithFetchFunc(scraping.FetchWithZyteProxy))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	result, err := scrapeai.Scrape(ctx, req)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}
	fmt.Printf("Zyte Proxy Result: %+v\n", result.Results)
}
