// Package scrapeai provides a public interface for web scraping with AI assistance.
// Version: v0.4.0
package scrapeai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samredway/scrapeai/gpt"
	"github.com/samredway/scrapeai/scraping"
)

const Version = "v0.4.0"

// ScrapeAiResult contains the results of a scraping operation.
type ScrapeAiResult struct {
	Url     string
	Results string
}

// Scrape performs a web scraping operation with AI assistance.
func Scrape(ctx context.Context, req *ScrapeAiRequest) (*ScrapeAiResult, error) {
	page, err := req.FetchFunc(ctx, req.Url)
	if err != nil {
		return nil, fmt.Errorf("fetching page: %w", err)
	}

	// scrape the page
	goqueryDoc, err := scraping.GoQueryDocFromBody(page)
	if err != nil {
		return nil, fmt.Errorf("creating goquery doc: %w", err)
	}
	strippedPage, err := scraping.StripNonTextTags(goqueryDoc)
	if err != nil {
		return nil, fmt.Errorf("stripping non-text tags: %w", err)
	}
	pageText, err := scraping.GetDocumentHTML(strippedPage)
	if err != nil {
		return nil, fmt.Errorf("getting document HTML: %w", err)
	}

	// get the results of search from GPT
	results, err := processWithGPT(req, pageText)
	if err != nil {
		return nil, fmt.Errorf("processing with GPT: %w", err)
	}

	return &ScrapeAiResult{
		Url:     req.Url,
		Results: results,
	}, nil
}

func processWithGPT(req *ScrapeAiRequest, pageText string) (string, error) {
	gptRequest := gpt.NewGptRequest(req.Prompt, pageText)
	if req.Schema != "" {
		gptRequest.SetSchema(req.Schema)
	}

	response, err := gpt.SendGptRequest(gptRequest)
	if err != nil {
		return "", err
	}

	// Get the raw response content
	content := response.Choices[0].Message.Content

	// Unmarshal into a generic structure to validate JSON
	var rawJSON any
	if err := json.Unmarshal([]byte(content), &rawJSON); err != nil {
		return "", fmt.Errorf("invalid JSON response: %w", err)
	}

	// If we get here, the JSON is valid
	return content, nil
}
