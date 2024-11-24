package gpt_test

import (
	"encoding/json"
	"strings"
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
		})
	}
}

func TestSchemaValidation(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		inError string
	}{
		{
			name:    "invalid json",
			schema:  `{"invalid-json"}`,
			inError: "invalid character '}'",
		},
		{
			name: "missing additional properties at top level",
			schema: `{
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
				"required": ["data"]
			}`,
			inError: "an object must contain the additionalProperties field",
		},
		{
			name: "missing additional properties at nested level",
			schema: `{
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
							"required": ["headline", "body"]
						}
					}
				},
				"additionalProperties": false,
				"required": ["data"]
			}`,
			inError: "an object must contain the additionalProperties field",
		},
		{
			name: "Missing a required properties array",
			schema: `{
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
							"required": ["body", "headline"],
							"additionalProperties": false
						}
					}
				},
				"additionalProperties": false
			}`,
			inError: "each value in properties must be in the required array",
		},
		{
			name: "Missing 'headline' field in required properties array",
			schema: `{
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
							"required": ["body"],
							"additionalProperties": false
						}
					}
				},
				"additionalProperties": false,
				"required": ["data"]
			}`,
			inError: "each value in properties must be in the required array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gpt.ValidateSchema(tt.schema)
			if err == nil {
				t.Errorf("Invalid schema did not cause an error!")
				return
			}
			if !strings.Contains(err.Error(), tt.inError) {
				t.Errorf("Expected error containing '%s' but got '%s'", tt.inError, err.Error())
			}
		})
	}
}
