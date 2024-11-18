package gpt

import (
	"encoding/json"
	"fmt"
)

// Validate the response schema that we provide to GPT for return
// data
func ValidateSchema(schema string) error {
	// check that the schema is valid json
	var obj any
	err := json.Unmarshal([]byte(schema), &obj)
	if err != nil {
		return fmt.Errorf("schema is not valid json: %w", err)
	}
	return nil
}
