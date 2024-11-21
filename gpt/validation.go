package gpt

import (
	"encoding/json"
	"fmt"
)

// Validate the response schema that we provide to GPT for return data
func ValidateSchema(schema string) error {
	var obj any
	err := json.Unmarshal([]byte(schema), &obj)
	if err != nil {
		return fmt.Errorf("schema is not valid json: %w", err)
	}
	return checkByType(obj.(map[string]any))
}

func checkByType(obj map[string]any) error {
	value, exists := obj["type"]
	if !exists {
		return checkNested(obj)
	}
	switch value {
	case "object":
		return checkObj(obj)
	case "array":
		return checkArr(obj)
	case "string":
		// TODO
		return nil
	default:
		return fmt.Errorf("Invalid value for key 'type' must be 'object' 'array' or 'string'")
	}
}

func checkNested(obj map[string]any) error {
	excludedKeys := map[string]bool{
		"additionalProperties": true,
		"type":                 true,
		"required":             true,
	}
	for key, value := range obj {
		if excludedKeys[key] {
			continue
		}
		if nestedMap, ok := value.(map[string]any); ok {
			return checkByType(nestedMap)
		} else {
			return fmt.Errorf("Expected value for key %s to be map[string]any, but got %T", key, value)
		}
	}
	return nil
}

func checkObj(obj map[string]any) error {
	value, exists := obj["additionalProperties"]
	if !exists || value != false {
		return fmt.Errorf("An object must contain the additionalProperties field and it must be false")
	}
	return checkNested(obj)
}

func checkArr(obj map[string]any) error {
	_, exists := obj["items"]
	if !exists {
		return fmt.Errorf("An array object must contain the key 'items'")
	}
	return checkNested(obj)
}
