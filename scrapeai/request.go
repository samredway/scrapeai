package scrapeai

import (
	"github.com/samredway/scrapeai/gpt"
	"github.com/samredway/scrapeai/scraping"
)

// functional options for ScrapeAiRequest
type Option func(*ScrapeAiRequest)

// FetchFunc is a function type for fetching a web page
// See scraping/utils/FetchFromChromedp for the default implementation
type FetchFunc func(url string) (string, error)

// default fetch function
var defaultFetchFunc = scraping.FetchFromChromedp

// Allows specifying a fetch function for collecting the web page the default
// is scraping.FetchFromChromedp
func WithFetchFunc(f FetchFunc) Option {
	return func(r *ScrapeAiRequest) {
		r.FetchFunc = f
	}
}

// Allows specifying a response schema for the return data. The default is a
// list of strings
func WithSchema(s string) Option {
	return func(r *ScrapeAiRequest) {
		r.Schema = s
	}
}

// ScrapeAiRequest represents the input for a scraping operation.
type ScrapeAiRequest struct {
	Url       string
	Prompt    string
	FetchFunc FetchFunc // Optional custom fetch function
	Schema    string    // Optional custom schema for the response
}

// Initialise a new ScrapeAiRequest object with options and sensible
// defaults
func NewScrapeAiRequest(url string, prompt string, options ...Option) (*ScrapeAiRequest, error) {
	req := &ScrapeAiRequest{Url: url, Prompt: prompt}
	for _, o := range options {
		o(req)
	}
	if req.FetchFunc == nil {
		req.FetchFunc = defaultFetchFunc
	}
	if req.Schema != "" {
		err := gpt.ValidateSchema(req.Schema)
		if err != nil {
			return nil, err
		}
	} else {
		req.Schema = gpt.DefaultSchemaTemplate
	}
	return req, nil
}
