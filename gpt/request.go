package gpt

import (
	"encoding/json"
	"fmt"
)

const DefaultSchemaTemplate = `{
    "type": "object",
    "properties": {
        "data": {
            "type": "array",
            "items": {"type": "string"}
        }
    },
    "additionalProperties": false,
    "required": ["data"]
}`

type GptRequest struct {
	Model          string         `json:"model"`
	Temperature    float64        `json:"temperature"`
	Seed           int            `json:"seed"`
	Messages       []GptMessage   `json:"messages"`
	ResponseFormat ResponseFormat `json:"response_format"`
}

type ResponseFormat struct {
	Type       string     `json:"type"`
	JSONSchema JsonSchema `json:"json_schema"`
}

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type JsonSchema struct {
	Name   string          `json:"name"`
	Strict bool            `json:"strict"`
	Schema json.RawMessage `json:"schema"`
}

var defaultRequest = GptRequest{
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
	request := defaultRequest
	request.Messages = []GptMessage{{Role: "user", Content: fmt.Sprintf(prompt, prompt, page)}}
	return &request
}

// DefaultSchema returns the default schema configuration which is a array of strings
func DefaultSchema() json.RawMessage {
	return json.RawMessage(DefaultSchemaTemplate)
}

// SetSchema sets a custom schema for the response from GPT.
// The schema can be pre-validated using ValidateSchema()
// before being set to avoid potential runtime errors.
func (r *GptRequest) SetSchema(schema string) {
	r.ResponseFormat.JSONSchema.Schema = json.RawMessage(schema)
}
