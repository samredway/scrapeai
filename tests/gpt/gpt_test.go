package gpt_test

import (
	"encoding/json"
	"testing"

	"github.com/samredway/scrapeai/gpt"
)

func TestGpt(t *testing.T) {
	tests := []struct {
		name   string
		schema string
		data   any
	}{
		{
			// Default schema is a slice of strings
			name:   "Default Schema",
			schema: "",
			data: &struct {
				Data []string `json:"data"`
			}{},
		},
		{
			// Test with a custom schema
			name: "Custom Schema",
			schema: `{
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
			}`,
			data: &struct {
				Data []struct {
					Headline string `json:"headline"`
					Body     string `json:"body"`
				} `json:"data"`
			}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := "Extract headline and body from this page"
			url := "https://example.com"

			request := gpt.NewGptRequest(prompt, url)
			if tt.schema != "" {
				request.SetSchema(tt.schema)
			}

			response, err := gpt.SendGptRequest(request)
			if err != nil {
				t.Errorf("Error sending GPT request: %v", err)
			}
			if response == nil {
				t.Errorf("No response from GPT")
				return
			}

			err = json.Unmarshal([]byte(response.Choices[0].Message.Content), tt.data)
			if err != nil {
				t.Errorf("Error unmarshalling response: %v", err)
			}

			t.Logf("Response: %+v", response.Choices[0].Message.Content)
		})
	}
}
