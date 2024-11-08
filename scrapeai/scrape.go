// Package scrapeai provides a public interface for web scraping with AI assistance.
// Version: v0.1.1
package scrapeai

import (
	"encoding/json"
	"fmt"

	"github.com/samredway/scrapeai/gpt"
	"github.com/samredway/scrapeai/scraping"
)

const Version = "v0.1.1"

// FetchFunc is a function type for fetching a web page
// See scraping/utils/FetchFromChromedp for the default implementation
type FetchFunc func(url string) (string, error)

// ScrapeAiRequest represents the input for a scraping operation.
type ScrapeAiRequest struct {
	Url       string
	Prompt    string
	FetchFunc FetchFunc // Optional custom fetch function
}

// ScrapeAiResult contains the results of a scraping operation.
type ScrapeAiResult struct {
	Url     string
	Results []string
}

// Scrape performs a web scraping operation with AI assistance.
func Scrape(req ScrapeAiRequest) (*ScrapeAiResult, error) {
	fetchFunc := req.FetchFunc
	if fetchFunc == nil {
		fetchFunc = scraping.FetchFromChromedp // Default fetch function
	}

	page, err := fetchFunc(req.Url)
	if err != nil {
		return nil, fmt.Errorf("fetching page: %w", err)
	}

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

	results, err := processWithGPT(req.Prompt, pageText)
	if err != nil {
		return nil, fmt.Errorf("processing with GPT: %w", err)
	}

	return &ScrapeAiResult{
		Url:     req.Url,
		Results: results,
	}, nil
}

func processWithGPT(prompt, pageText string) ([]string, error) {
	gptRequest := gpt.NewGptRequest(prompt, pageText)
	response, err := gpt.SendGPTRequest(&gptRequest)
	if err != nil {
		return nil, err
	}

	var jsonResponse struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &jsonResponse); err != nil {
		return nil, fmt.Errorf("unmarshaling JSON response: %w", err)
	}

	return jsonResponse.Data, nil
}
