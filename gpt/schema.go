package gpt

// SchemaObject is the definition of the return object for GPT requests
// It is used to define the structure of the response from GPT
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

// NewSchemaObject allows users to define a custom schema for the response
func NewSchemaObject() *SchemaObject {
	// TODO: allow this to be a customer object and move a way of creating
	// this default schema into the DefaultSchema function
	return &SchemaObject{
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
	}
}

// DefaultSchema returns the default schema configuration which is a array of strings
func DefaultSchema() *SchemaObject {
	return NewSchemaObject()
}
