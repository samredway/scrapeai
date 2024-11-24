package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/samredway/scrapeai/scrapeai"
	"github.com/samredway/scrapeai/scraping"
)

func main() {
	url := "https://example.com"

	// Example: Custom Schema
	fmt.Println("=== Custom Schema Example ===")
	type ReturnData struct {
		Headline string `json:"headline"`
		Body     string `json:"body"`
	}

	schema := `
		{
			"type": "object",
			"properties": {
				"headline": {"type": "string"},
				"body": {"type": "string"}
			},
			"additionalProperties": false,
			"required": ["headline", "body"]
		}
	`

	req, err := scrapeai.NewScrapeAiRequest(url,
		"Extract the headline and body content",
		scrapeai.WithFetchFunc(scraping.Fetch),
		scrapeai.WithSchema(schema),
	)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	result, err := scrapeai.Scrape(req)
	if err != nil {
		log.Fatalf("Error scraping with AI: %v", err)
	}

	var content ReturnData
	err = json.Unmarshal([]byte(result.Results), &content)
	if err != nil {
		log.Fatalf("Failed to convert result to ReturnData: %v", err)
	}

	fmt.Printf("Custom Schema Result:\nHeadline: %s\nBody: %s\n",
		content.Headline,
		content.Body,
	)
}
