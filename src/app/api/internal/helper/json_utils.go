package helper

import (
	"fmt"
	"strconv"
	"strings"
)

func EscapeJSON(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
	)
	return replacer.Replace(s)
}

func ExtractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}
	return "{}"
}

func ExtractByPath(data interface{}, path string) string {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if current == nil {
			return ""
		}

		// Handle array indexing (e.g., candidates[0])
		if strings.Contains(part, "[") {
			idxStart := strings.Index(part, "[")
			idxEnd := strings.Index(part, "]")
			arrName := part[:idxStart]

			var idx int
			fmt.Sscanf(part[idxStart+1:idxEnd], "%d", &idx)

			m, ok := current.(map[string]interface{})
			if !ok {
				return ""
			}

			arr, ok := m[arrName].([]interface{})
			if !ok || idx >= len(arr) {
				return ""
			}

			current = arr[idx]
			continue
		}

		// Handle numeric segment
		if idx, err := strconv.Atoi(part); err == nil {
			arr, ok := current.([]interface{})
			if !ok || idx >= len(arr) {
				return ""
			}
			current = arr[idx]
			continue
		}

		// Handle map access
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current = m[part]
	}

	if s, ok := current.(string); ok {
		return s
	}
	return ""
}

func ReplaceTemplate(tpl string, vars map[string]string) string {
	for k, v := range vars {
		tpl = strings.ReplaceAll(tpl, k, v)
	}
	return tpl
}
