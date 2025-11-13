package gpt_test

import (
	"testing"

	"github.com/samredway/scrapeai/gpt"
)

func TestNewGptRequest(t *testing.T) {
	tests := []struct {
		prompt   string
		page     string
		expected string
	}{
		{
			prompt:   "Extract the main headline",
			page:     "<html><body>Example Domain</body></html>",
			expected: "Extract the main headline\n\n<html><body>Example Domain</body></html>",
		},
		{
			prompt:   "Extract prices (discount is 50% off)",
			page:     "<html><body>$100</body></html>",
			expected: "Extract prices (discount is 50% off)\n\n<html><body>$100</body></html>",
		},
	}

	for _, tt := range tests {
		request := gpt.NewGptRequest(tt.prompt, tt.page)

		if len(request.Messages) != 1 {
			t.Fatalf("Expected exactly one message, got %d", len(request.Messages))
		}

		content := request.Messages[0].Content
		if content != tt.expected {
			t.Errorf("Expected message content to be exactly %q, but got %q", tt.expected, content)
		}

	}

}
