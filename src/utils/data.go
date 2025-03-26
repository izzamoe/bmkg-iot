package utils

import (
	"log"
	"strings"

	"github.com/goccy/go-json"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// structToMap converts a struct to a map[string]interface{} for flexible data handling
func StructToMap(obj interface{}) map[string]interface{} {
	// Convert to JSON and back to map for flexibility
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Error marshaling struct to JSON: %v", err)
		return map[string]interface{}{}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		log.Printf("Error unmarshaling JSON to map: %v", err)
		return map[string]interface{}{}
	}

	return result
}

// MapToStruct converts a map[string]interface{} to a struct using generics
func MapToStruct[T any](data map[string]interface{}) (T, error) {
	var result T

	// Convert map to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	// Convert JSON to struct
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, err
	}

	return result, nil
}

// IsZeroOrNil checks if a value is the zero value for its type or nil
func IsZeroOrNil(v interface{}) bool {
	if v == nil {
		return true
	}

	switch v := v.(type) {
	case int, int8, int16, int32, int64:
		return v == 0
	case uint, uint8, uint16, uint32, uint64:
		return v == 0
	case float32, float64:
		return v == 0
	case bool:
		return v == false
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		// For complex types, convert to JSON to see if it's empty
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			log.Printf("Error marshaling in IsZeroOrNil: %v", err)
			return false
		}
		return string(jsonBytes) == "{}" || string(jsonBytes) == "[]" || string(jsonBytes) == "null"
	}
}

// DeepMerge merges two maps recursively, with values from 'src' overriding 'dst'
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	for key, srcVal := range src {
		dstVal, exists := dst[key]

		if exists {
			// If both values are maps, merge them
			srcMap, srcIsMap := srcVal.(map[string]interface{})
			dstMap, dstIsMap := dstVal.(map[string]interface{})

			if srcIsMap && dstIsMap {
				dst[key] = DeepMerge(dstMap, srcMap)
				continue
			}
		}

		// Otherwise just replace the value
		dst[key] = srcVal
	}

	return dst
}

// FilterMap creates a new map containing only the keys specified in the filter
func FilterMap(data map[string]interface{}, keys []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range keys {
		if val, exists := data[key]; exists {
			result[key] = val
		}
	}

	return result
}

// SafeGet retrieves a nested value from a map using a dot-notation path with type safety
func SafeGet[T any](data map[string]interface{}, path string, defaultValue T) T {
	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys[:len(keys)-1] {
		val, exists := current[key]
		if !exists {
			return defaultValue
		}

		// Check if we can continue traversing
		nextMap, ok := val.(map[string]interface{})
		if !ok {
			return defaultValue
		}
		current = nextMap
	}

	// Get the final value
	lastKey := keys[len(keys)-1]
	if val, exists := current[lastKey]; exists {
		if result, ok := val.(T); ok {
			return result
		}
	}

	return defaultValue
}

// OmitFields creates a new map without the specified keys
func OmitFields(data map[string]interface{}, keys []string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range data {
		excluded := false
		for _, excludeKey := range keys {
			if k == excludeKey {
				excluded = true
				break
			}
		}

		if !excluded {
			result[k] = v
		}
	}

	return result
}

// FlattenMap converts a nested map to a flat map with dot notation keys
func FlattenMap(data map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		if nestedMap, ok := v.(map[string]interface{}); ok {
			// Recursively flatten nested maps
			flatMap := FlattenMap(nestedMap, key)
			for nestedKey, nestedVal := range flatMap {
				result[nestedKey] = nestedVal
			}
		} else {
			result[key] = v
		}
	}

	return result
}
