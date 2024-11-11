package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const openaiApiUrl = "https://api.openai.com/v1/chat/completions"

// SendGptRequest sends a GPT request to the OpenAI API and returns a GptResponse
// TODO: Handle size limits (chunking strategy)
func SendGptRequest(config *GptRequest) (*GptResponse, error) {
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error with request to GPT API unable to read body of response")
		}
		return nil, fmt.Errorf("GPT API request failed with status code: %d and message %s", resp.StatusCode, body)
	}

	var gptResponse GptResponse
	if err := json.NewDecoder(resp.Body).Decode(&gptResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &gptResponse, nil
}
