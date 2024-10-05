// Package scrapeai contains the Public interface for the scrapeai package
package scrapeai

import (
	"strings"

	"github.com/samredway/scrapeai/scraping"
)

type ScrapeAiRequest struct {
	Url    string
	Prompt string
}

type ScrapeAiResult struct {
	Url  string
	Text []string
}

func Scrape(req ScrapeAiRequest) (ScrapeAiResult, error) {
	page, err := scraping.FetchFromChomedp(req.Url)
	if err != nil {
		return ScrapeAiResult{}, err
	}
	goqueryDoc, err := scraping.GoQueryDocFromBody(page)
	if err != nil {
		return ScrapeAiResult{}, err
	}
	strippedPage, err := scraping.StripNonTextTags(goqueryDoc)
	if err != nil {
		return ScrapeAiResult{}, err
	}
	pageText, err := scraping.GetDocumentHTML(strippedPage)
	if err != nil {
		return ScrapeAiResult{}, err
	}
	gptRequest := newGptRequest(req.Prompt, pageText)
	response, err := generateText(&gptRequest)
	if err != nil {
		return ScrapeAiResult{}, err
	}
	return ScrapeAiResult{Url: req.Url, Text: strings.Split(response, gptResponseSeparator)}, nil
}
