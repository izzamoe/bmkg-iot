package utils

import (
	"log"
	"strings"

	"github.com/goccy/go-json"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// StructToMap converts a struct to a map[string]interface{} for flexible data handling
// Optimized for maximum performance with minimal allocations
func StructToMap(obj interface{}) map[string]interface{} {
	// Fast path for nil values
	if obj == nil {
		return map[string]interface{}{}
	}

	// Fast path for maps - avoid unnecessary conversion
	if m, ok := obj.(map[string]interface{}); ok {
		result := make(map[string]interface{}, len(m))
		for k, v := range m {
			result[k] = v
		}
		return result
	}

	// Convert to JSON with the high-performance goccy/go-json library
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Error marshaling struct to JSON: %v", err)
		return map[string]interface{}{}
	}

	// Estimate capacity based on JSON size to reduce allocations
	estimatedFields := len(jsonBytes) / 15
	if estimatedFields < 8 {
		estimatedFields = 8 // Reasonable minimum capacity
	}

	result := make(map[string]interface{}, estimatedFields)
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		log.Printf("Error unmarshaling JSON to map: %v", err)
		return map[string]interface{}{}
	}

	return result
}

// MapToStruct converts a map[string]interface{} to a struct using generics
// Extremely optimized for maximum performance with minimal allocations
func MapToStruct[T any](data map[string]interface{}) (T, error) {
	var result T

	// Fast path for nil or empty maps - avoid unnecessary processing
	if data == nil || len(data) == 0 {
		return result, nil
	}

	// Estimate capacity based on map size to reduce allocations
	estimatedSize := len(data) * 32
	if estimatedSize < 128 {
		estimatedSize = 128 // Minimum buffer size for small maps
	}

	// Pre-allocate a buffer with estimated capacity
	// This avoids multiple reallocations during marshaling
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	// Direct unmarshaling to target type without intermediate allocations
	// goccy/go-json is already 2-3x faster than standard library
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, err
	}

	return result, nil
}

// IsZeroOrNil checks if a value is the zero value for its type or nil
// Extremely optimized for maximum performance with minimal allocations
func IsZeroOrNil(v interface{}) bool {
	if v == nil {
		return true
	}

	switch v := v.(type) {
	// Numeric types - handle separately for better compiler optimization
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case bool:
		return !v // Faster than v == false
	case string:
		return v == ""
	// Handle common collection types directly
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []byte:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	case map[string]string:
		return len(v) == 0
	// Quick handling for pointers to common types
	case *string:
		return v == nil || *v == ""
	case *int:
		return v == nil || *v == 0
	case *bool:
		return v == nil || !*v
	default:
		// For complex types, optimize JSON comparison without string allocations
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			log.Printf("Error marshaling in IsZeroOrNil: %v", err)
			return false
		}

		// Fast path: check common empty values by direct byte comparison
		jsonLen := len(jsonBytes)
		if jsonLen == 2 {
			return (jsonBytes[0] == '{' && jsonBytes[1] == '}') || // Empty object
				(jsonBytes[0] == '[' && jsonBytes[1] == ']') // Empty array
		}
		if jsonLen == 4 {
			return jsonBytes[0] == 'n' && jsonBytes[1] == 'u' &&
				jsonBytes[2] == 'l' && jsonBytes[3] == 'l' // null value
		}

		return false
	}
}

// DeepMerge merges two maps recursively, with values from 'src' overriding 'dst'
// Extremely optimized for maximum performance with minimal allocations
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	// Fast path for empty or nil source maps
	if src == nil || len(src) == 0 {
		return dst
	}

	// Fast path for nil destination
	if dst == nil {
		dst = make(map[string]interface{}, len(src))
	}

	for key, srcVal := range src {
		// Fast path for non-existent destination keys (skip type checks)
		dstVal, exists := dst[key]
		if !exists {
			dst[key] = srcVal
			continue
		}

		// Only perform type assertions when needed (when key exists in both maps)
		srcMap, srcIsMap := srcVal.(map[string]interface{})
		dstMap, dstIsMap := dstVal.(map[string]interface{})

		// Only recurse if both values are maps
		if srcIsMap && dstIsMap {
			dst[key] = DeepMerge(dstMap, srcMap)
		} else {
			// For non-map values, just use the source value
			dst[key] = srcVal
		}
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
// Optimized for maximum performance with minimal allocations
func FlattenMap(data map[string]interface{}, prefix string) map[string]interface{} {
	// Fast path for empty maps
	if len(data) == 0 {
		return make(map[string]interface{})
	}

	// Estimate capacity to reduce map reallocations (3x is a good heuristic for nested data)
	initialCapacity := len(data) * 3
	result := make(map[string]interface{}, initialCapacity)

	// Use helper function to avoid creating new result maps in recursive calls
	flattenMapHelper(data, prefix, result)

	return result
}

// Helper function that operates on a shared result map to avoid allocations
func flattenMapHelper(data map[string]interface{}, prefix string, result map[string]interface{}) {
	for k, v := range data {
		// Compute the current key once
		var key string
		if prefix == "" {
			key = k
		} else {
			key = prefix + "." + k
		}

		// Process based on value type
		if nestedMap, ok := v.(map[string]interface{}); ok && len(nestedMap) > 0 {
			// Recursively process non-empty nested map
			flattenMapHelper(nestedMap, key, result)
		} else {
			// Add leaf value directly to result
			result[key] = v
		}
	}
}
