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

func processWithGPT(req *ScrapeAiRequest, pageText string) ([]string, error) {
	gptRequest := gpt.NewGptRequest(req.Prompt, pageText)
	if req.Schema != "" {
		gptRequest.SetSchema(req.Schema)
	}

	response, err := gpt.SendGptRequest(gptRequest)
	if err != nil {
		return nil, err
	}

	// this will return a slice of the raw response strings that have been returned
	// the caller will need to unmarshal the response to match the schema that is
	// passed in
	var jsonResponse struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &jsonResponse); err != nil {
		return nil, fmt.Errorf("unmarshaling JSON response: %w", err)
	}

	return jsonResponse.Data, nil
}
