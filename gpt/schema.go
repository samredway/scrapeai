package gpt

import "encoding/json"

const DefaultSchemaTemplate = `{"type": "array", "items": {"type": "string"}}`

// SchemaObject is the definition of the return object for GPT requests
// It is used to define the structure of the response from GPT
type SchemaObject struct {
	Type                 string             `json:"type"`
	Properties           SchemaDataProperty `json:"properties"`
	AdditionalProperties bool               `json:"additionalProperties"`
	Required             []string           `json:"required"`
}

type SchemaDataProperty struct {
	Data json.RawMessage `json:"data"`
}

// NewSchemaObject allows users to define a custom schema for the response
func NewSchemaObject(schema json.RawMessage) *SchemaObject {
	return &SchemaObject{
		Type:                 "object",
		Properties:           SchemaDataProperty{Data: schema},
		AdditionalProperties: false,
		Required:             []string{"data"},
	}
}

// DefaultSchema returns the default schema configuration which is a array of strings
func DefaultSchema() *SchemaObject {
	defaultSchema := json.RawMessage(DefaultSchemaTemplate)
	return NewSchemaObject(defaultSchema)
}
