// Package scrapeai provides a public interface for web scraping with AI assistance.
// Version: v0.1.2
package scrapeai

import (
	"encoding/json"
	"fmt"

	"github.com/samredway/scrapeai/gpt"
	"github.com/samredway/scrapeai/scraping"
)

const Version = "v0.1.2"

// ScrapeAiResult contains the results of a scraping operation.
type ScrapeAiResult struct {
	Url     string
	Results []string
}

// Scrape performs a web scraping operation with AI assistance.
func Scrape(req *ScrapeAiRequest) (*ScrapeAiResult, error) {
	page, err := req.FetchFunc(req.Url)
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
	response, err := gpt.SendGptRequest(gptRequest)
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
