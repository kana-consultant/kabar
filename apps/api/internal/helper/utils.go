package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Helper functions
func getMapKeys(m map[string]interface{}) []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}

func isSensitiveField(fieldName string) bool {
	return sensitiveFields[fieldName]
}

func removeSensitiveFields(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range data {
		if isSensitiveField(k) {
			result[k] = "[REDACTED]"
		} else if nestedMap, ok := v.(map[string]interface{}); ok {
			result[k] = removeSensitiveFields(nestedMap)
		} else {
			result[k] = v
		}
	}
	return result
}

func redactSensitiveFields(body map[string]interface{}) map[string]interface{} {
	return removeSensitiveFields(body)
}

func stripHTML(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	clean := re.ReplaceAllString(input, "")
	clean = strings.ReplaceAll(clean, "\n", " ")
	clean = strings.ReplaceAll(clean, "\r", " ")
	clean = strings.Join(strings.Fields(clean), " ")
	return clean
}

func getStringValue(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			return strVal
		}
	}
	return defaultValue
}

func resolveTemplate(v string, data map[string]string) string {
	result := v
	for k, val := range data {
		result = strings.ReplaceAll(result, "{{"+k+"}}", val)
	}
	return result
}

func printJSON(title string, data interface{}) {
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s:\n%s\n", title, string(jsonBytes))
}
