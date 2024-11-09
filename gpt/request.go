package gpt

import "fmt"

type GptRequest struct {
	Model          string         `json:"model"`
	Temperature    float64        `json:"temperature"`
	Seed           int            `json:"seed"`
	Messages       []GptMessage   `json:"messages"`
	ResponseFormat ResponseFormat `json:"response_format"`
}

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type       string     `json:"type"`
	JSONSchema JsonSchema `json:"json_schema"`
}

type JsonSchema struct {
	Name   string        `json:"name"`
	Strict bool          `json:"strict"`
	Schema *SchemaObject `json:"schema"`
}

var defaultConfig = GptRequest{
	Model:       "gpt-4o-mini",
	Temperature: 0.0,
	Seed:        42,
	ResponseFormat: ResponseFormat{
		Type: "json_schema",
		JSONSchema: JsonSchema{
			Name:   "scrape_result",
			Strict: true,
			Schema: DefaultSchema(),
		},
	},
}

// NewGptRequest creates a new GPT request with the given prompt and page
func NewGptRequest(prompt, page string) *GptRequest {
	config := defaultConfig
	config.Messages = []GptMessage{{Role: "user", Content: fmt.Sprintf(prompt, prompt, page)}}
	return &config
}
