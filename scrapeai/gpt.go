package scrapeai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	openaiApiUrl         = "https://api.openai.com/v1/chat/completions"
	gptResponseSeparator = ";;;"
	prompt               = `
	You are a helpful assistant that can scrape the web.
	You are given a URL and a prompt.
	You need to scrape the web page and return the text as
	outlined in the prompt.

	There may be multiple instances of the text you need to
	extract. If so, please return all of them. Please return
	them using the following separator: ';;;'.

	An example prompt is:
	"Extract the job titles from the page."

	An example response is:
	"Software Engineer;;;Product Manager;;;"

	Here is your prompt:
	%s

	The following text is the raw html of the page you must scrape:
	%s

	Please return the text as outlined in the prompt.
	`
)

var openaiApiKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	openaiApiKey = os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set in the environment")
	}
}

type gptRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Seed        int       `json:"seed"`
	Messages    []gptMessage `json:"messages"`
}

type gptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var defaultConfig = gptRequest{
	Model:       "gpt-4o-mini",
	MaxTokens:   1024,
	Temperature: 0.0,
	Seed:        42,
}

func newGptRequest(prompt string, page string) gptRequest {
	config := defaultConfig
	config.Messages = []gptMessage{{Role: "user", Content: fmt.Sprintf(prompt, prompt, page)}}
	return config
}

// TODO return a GPT result object with properly parsed data
func generateText(config *gptRequest) (string, error) {
	body, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openaiApiUrl, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}
