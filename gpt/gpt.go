package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const openaiApiUrl = "https://api.openai.com/v1/chat/completions"

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

func NewGptRequest(prompt, page string) GptRequest {
	config := defaultConfig
	config.Messages = []GptMessage{{Role: "user", Content: fmt.Sprintf(prompt, prompt, page)}}
	return config
}

// TODO: Return a GPT result object with properly parsed data
// TODO: Handle size limits (chunking strategy)
func SendGPTRequest(config *GptRequest) (*GptResponse, error) {
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
