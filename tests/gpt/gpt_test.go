package gpt_test

import (
	"fmt"
	"testing"

	"github.com/samredway/scrapeai/gpt"
)

func TestGpt(t *testing.T) {
	prompt := "Extract headline and body from this page"
	url := "https://example.com"

	// response_schema := "{\"headline\": \"string\", \"body\": \"string\"}"
	request := gpt.NewGptRequest(prompt, url)
	response, err := gpt.SendGptRequest(request)
	if err != nil {
		t.Errorf("Error sending GPT request: %v", err)
	}

	fmt.Println(response)
}
