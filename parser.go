package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ParseDynamicFlags parses command-line flags into a map
func ParseDynamicFlags(args []string) map[string]interface{} {
	flags := make(map[string]interface{})
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(arg, "--")
			var value string
			if strings.Contains(key, "=") {
				parts := strings.SplitN(key, "=", 2)
				key = parts[0]
				value = parts[1]
			} else {
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
					i++
					value = args[i]
				} else {
					value = ""
				}
			}
			flags[key] = parseValue(value)
		}
		i++
	}
	return flags
}

// parseValue attempts to parse a string value into an appropriate type
func parseValue(value string) interface{} {
	// Try parsing as int
	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		return intVal
	}
	// Try parsing as bool
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}
	// Try parsing as float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	// Try parsing as JSON (array or object)
	trimmedValue := strings.TrimSpace(value)
	if strings.HasPrefix(trimmedValue, "[") || strings.HasPrefix(trimmedValue, "{") {
		var jsonData interface{}
		err := json.Unmarshal([]byte(trimmedValue), &jsonData)
		if err == nil {
			return jsonData
		}
	}
	// Return as string
	return value
}
