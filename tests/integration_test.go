package integration_test

import (
	"strings"
	"testing"

	"github.com/samredway/scrapeai/scrapeai"
	"github.com/samredway/scrapeai/scraping"
)

const exampleUrl = "https://example.com"

func TestScrapeIntegration(t *testing.T) {

	tests := []struct {
		name           string
		prompt         string
		expectedPart   string
		unexpectedPart string
	}{
		{
			name:           "Extract main headline",
			prompt:         "Extract the main headline",
			expectedPart:   "Example Domain",
			unexpectedPart: "More information",
		},
		{
			name:           "Extract main body",
			prompt:         "Extract the main body of the page, exclude the main headline",
			expectedPart:   "This domain is for use in illustrative examples in documents",
			unexpectedPart: "Example Domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := scrapeai.NewScrapeAiRequest(exampleUrl, tt.prompt, scrapeai.WithFetchFunc(scraping.Fetch))
			result, err := scrapeai.Scrape(req)
			if err != nil {
				t.Fatalf("Error scraping with AI: %v", err)
			}

			if len(result.Results) == 0 {
				t.Fatalf("No results returned")
			}

			if !strings.Contains(result.Results[0], tt.expectedPart) {
				t.Errorf("Expected result to contain '%s', but it didn't. Got: %s", tt.expectedPart, result.Results[0])
			}

			if strings.Contains(result.Results[0], tt.unexpectedPart) {
				t.Errorf("Expected result not to contain '%s', but it did. Got: %s", tt.unexpectedPart, result.Results[0])
			}

			if result.Url != exampleUrl {
				t.Errorf("Expected url to be '%s', but got '%s'", exampleUrl, result.Url)
			}
		})
	}
}

func TestCustomSchema(t *testing.T) {
	req := scrapeai.NewScrapeAiRequest(
		exampleUrl,
		"Extract the headline and the body and return them in the specified data object",
		scrapeai.WithFetchFunc(scraping.Fetch),
		scrapeai.WithSchema(`{"type": "object", "headline": {"type": string}, "body": {"type": string}}`),
	)
	result, err := scrapeai.Scrape(req)
	if err != nil {
		t.Fatalf("Error scraping with AI: %v", err)
	}

	if len(result.Results) == 0 {
		t.Fatalf("No results returned")
	}

	println("Result:", result.Results)

}
