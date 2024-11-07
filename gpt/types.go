package gpt

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
