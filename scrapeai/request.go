package scrapeai

import "github.com/samredway/scrapeai/scraping"

// FetchFunc is a function type for fetching a web page
// See scraping/utils/FetchFromChromedp for the default implementation
type FetchFunc func(url string) (string, error)

// ScrapeAiRequest represents the input for a scraping operation.
type ScrapeAiRequest struct {
	Url       string
	Prompt    string
	FetchFunc FetchFunc // Optional custom fetch function
	Schema    string    // Optional custom schema for the response
}

// functional options for Scrape FetchFunc
type Option func(*ScrapeAiRequest)

// Allows specifying a fetch function for collecting the web page the default
// is scraping.FetchFromChromedp
func WithFetchFunc(f FetchFunc) Option {
	return func(r *ScrapeAiRequest) {
		r.FetchFunc = f
	}
}

var defaultFetchFunc = scraping.FetchFromChromedp

// Allows specifying a response schema for the return data. The default is a
// list of strings
func WithSchema(s string) Option {
	return func(r *ScrapeAiRequest) {
		r.Schema = s
	}
}

const defaultSchema = `{"type": "array", "items": {"type": "string"}}`

// Initialise a new ScrapeAiRequest object with options and sensible
// defaults
func NewScrapeAiRequest(url string, prompt string, options ...Option) *ScrapeAiRequest {
	req := &ScrapeAiRequest{Url: url, Prompt: prompt}
	for _, o := range options {
		o(req)
	}
	if req.FetchFunc == nil {
		req.FetchFunc = defaultFetchFunc
	}
	if req.Schema == "" {
		req.Schema = defaultSchema
	}
	return req
}
