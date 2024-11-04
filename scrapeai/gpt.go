package scrapeai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const openaiApiUrl = "https://api.openai.com/v1/chat/completions"

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
	JSONSchema JSONSchema `json:"json_schema"`
}

type JSONSchema struct {
	Name   string       `json:"name"`
	Strict bool         `json:"strict"`
	Schema SchemaObject `json:"schema"`
}

type SchemaObject struct {
	Type                 string             `json:"type"`
	Properties           SchemaDataProperty `json:"properties"`
	AdditionalProperties bool               `json:"additionalProperties"`
	Required             []string           `json:"required"`
}

type SchemaDataProperty struct {
	Data SchemaDataArray `json:"data"`
}

type SchemaDataArray struct {
	Type  string `json:"type"`
	Items struct {
		Type string `json:"type"`
	} `json:"items"`
}

type GptResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choice  `json:"choices"`
	Usage   UsageInfo `json:"usage"`
}

type Choice struct {
	Index        int        `json:"index"`
	Message      GptMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

var defaultConfig = GptRequest{
	Model:       "gpt-4o-mini",
	Temperature: 0.0,
	Seed:        42,
	ResponseFormat: ResponseFormat{
		Type: "json_schema",
		JSONSchema: JSONSchema{
			Name:   "scrape_result",
			Strict: true,
			Schema: SchemaObject{
				Type: "object",
				Properties: SchemaDataProperty{
					Data: SchemaDataArray{
						Type: "array",
						Items: struct {
							Type string `json:"type"`
						}{
							Type: "string",
						},
					},
				},
				AdditionalProperties: false,
				Required:             []string{"data"},
			},
		},
	},
}

func newGptRequest(prompt, page string) GptRequest {
	config := defaultConfig
	config.Messages = []GptMessage{{Role: "user", Content: fmt.Sprintf(prompt, prompt, page)}}
	return config
}

// TODO: Return a GPT result object with properly parsed data
// TODO: Handle size limits (chunking strategy)
func sendGPTRequest(config *GptRequest) (*GptResponse, error) {
	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set in the environment")
	}

	body, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", openaiApiUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var gptResponse GptResponse
	if err := json.NewDecoder(resp.Body).Decode(&gptResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &gptResponse, nil
}
