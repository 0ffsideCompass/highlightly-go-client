package client

import (
	"fmt"
)

// asInt converts an integer-like value (an int, or a whole-number float64 as
// produced by JSON decoding) to an int.
func asInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		if v == float64(int64(v)) {
			return int(v), nil
		}
		return 0, fmt.Errorf("value %v is not an integer", v)
	default:
		return 0, fmt.Errorf("value %v is not an integer", value)
	}
}

// requireIntParams ensures that every given key is present in params and holds
// an integer value (an int, or a whole-number float64 as produced by JSON
// decoding).
func requireIntParams(params map[string]any, keys ...string) error {
	for _, key := range keys {
		value, ok := params[key]
		if !ok {
			return fmt.Errorf("'%s' is required", key)
		}
		if _, err := asInt(value); err != nil {
			return fmt.Errorf("'%s' must be an integer", key)
		}
	}
	return nil
}

// requireStringParams ensures that every given key is present in params and
// holds a non-empty string value.
func requireStringParams(params map[string]any, keys ...string) error {
	for _, key := range keys {
		value, ok := params[key]
		if !ok {
			return fmt.Errorf("'%s' is required", key)
		}
		s, ok := value.(string)
		if !ok || s == "" {
			return fmt.Errorf("'%s' must be a non-empty string", key)
		}
	}
	return nil
}
