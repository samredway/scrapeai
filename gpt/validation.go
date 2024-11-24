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
		return fmt.Errorf("invalid value for key 'type' must be 'object' 'array' or 'string'")
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
			return fmt.Errorf("expected value for key %s to be map[string]any, but got %T", key, value)
		}
	}
	return nil
}

func checkObj(obj map[string]any) error {
	adProp, apExists := obj["additionalProperties"]
	if !apExists || adProp != false {
		return fmt.Errorf("an object must contain the additionalProperties field and it must be false")
	}
	props, pExists := obj["properties"]
	if !pExists {
		return fmt.Errorf("an object must contain the 'properties' field")
	}
	propsMap, ok := props.(map[string]any)
	if !ok {
		return fmt.Errorf("properties must be a valid object")
	}
	propKeys := make([]string, 0, len(propsMap))
	for key := range propsMap {
		propKeys = append(propKeys, key)
	}
	required, rExists := obj["required"]
	if !rExists && len(propKeys) > 0 {
		return fmt.Errorf("each value in properties must be in the required array")
	}
	// Verify all properties are required
	if rExists {
		requiredList, ok := required.([]any)
		if !ok {
			return fmt.Errorf("required must be an array")
		}
		if len(requiredList) != len(propKeys) {
			return fmt.Errorf("each value in properties must be in the required array")
		}
		requiredSet := make(map[string]bool)
		for _, r := range requiredList {
			if rStr, ok := r.(string); ok {
				requiredSet[rStr] = true
			}
		}
		for _, key := range propKeys {
			if !requiredSet[key] {
				return fmt.Errorf("property %s must be required", key)
			}
		}
	}
	return checkNested(obj)
}

func checkArr(obj map[string]any) error {
	_, exists := obj["items"]
	if !exists {
		return fmt.Errorf("an array object must contain the key 'items'")
	}
	return checkNested(obj)
}
