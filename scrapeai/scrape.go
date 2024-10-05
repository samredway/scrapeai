package scrapeai

import (
	"strings"
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
	page, err := FetchFromChomedp(req.Url)
	if err != nil {
		return ScrapeAiResult{}, err
	}
	gptRequest := newGptRequest(req.Prompt, page)
	response, err := generateText(&gptRequest)
	if err != nil {
		return ScrapeAiResult{}, err
	}
	return ScrapeAiResult{Url: req.Url, Text: strings.Split(response, ";;;")}, nil
}
