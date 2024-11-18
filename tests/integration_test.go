package integration_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/samredway/scrapeai/scrapeai"
	"github.com/samredway/scrapeai/scraping"
)

const exampleUrl = "https://example.com"

func TestScrapeDefaultSchema(t *testing.T) {
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

			data := result.Results

			t.Logf("Data: %s", data)

			type defaultSchema struct {
				Data []string `json:"data"`
			}

			var response defaultSchema
			err = json.Unmarshal([]byte(data), &response)
			if err != nil {
				t.Fatalf("Error unmarshalling JSON response: %v", err)
			}

			if len(response.Data) == 0 {
				t.Fatalf("No results returned")
			}

			if !strings.Contains(response.Data[0], tt.expectedPart) {
				t.Errorf("Expected result to contain '%s', but it didn't. Got: %s", tt.expectedPart, response.Data[0])
			}

			if strings.Contains(response.Data[0], tt.unexpectedPart) {
				t.Errorf("Expected result not to contain '%s', but it did. Got: %s", tt.unexpectedPart, response.Data[0])
			}

			if result.Url != exampleUrl {
				t.Errorf("Expected url to be '%s', but got '%s'", exampleUrl, result.Url)
			}
		})
	}
}

func TestScrapeCustomSchema(t *testing.T) {
	test_schema := `{
		"type": "object",
		"properties": {
			"data": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"headline": {"type": "string"},
						"body": {"type": "string"}
					},
					"additionalProperties": false,
					"required": ["headline", "body"]
				}
			}
		},
		"additionalProperties": false,
		"required": ["data"]
	}`
	req := scrapeai.NewScrapeAiRequest(
		exampleUrl,
		"Extract the headline and the body and return them in the specified data object",
		scrapeai.WithFetchFunc(scraping.Fetch),
		scrapeai.WithSchema(test_schema),
	)
	result, err := scrapeai.Scrape(req)
	if err != nil {
		t.Fatalf("Error scraping with AI: %v", err)
	}

	data := result.Results

	// Unmarshal the response to the expected schema will validate the response
	var jsonResponse struct {
		Data []map[string]string `json:"data"`
	}
	err = json.Unmarshal([]byte(data), &jsonResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON response: %v", err)
	}
}

func TestScrapeErrors(t *testing.T) {
	t.Run("invalid URL", func(t *testing.T) {
		req := scrapeai.NewScrapeAiRequest("not-a-url", "Extract headline")
		_, err := scrapeai.Scrape(req)
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})
}
