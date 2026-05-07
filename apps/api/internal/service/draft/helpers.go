package draft

import (
	"encoding/json"
	"fmt"
	"seo-backend/internal/models"
	"strings"
	"time"
)

func nullIfEmpty(id string) interface{} {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	return id
}

func prepareUpdateData(updates map[string]interface{}) map[string]interface{} {
	fieldMap := map[string]string{
		"title":          "title",
		"topic":          "topic",
		"article":        "article",
		"imageUrl":       "image_url",
		"imagePrompt":    "image_prompt",
		"status":         "status",
		"scheduledFor":   "scheduled_for",
		"targetProducts": "target_products",
		"hasImage":       "has_image",
	}

	data := make(map[string]interface{})
	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok {
			if key == "targetProducts" {
				jsonValue, _ := json.Marshal(value)
				data[dbField] = jsonValue
			} else if key == "scheduledFor" && value != nil {
				if scheduledStr, ok := value.(string); ok && scheduledStr != "" {
					if !strings.HasPrefix(scheduledStr, "daily:") {
						if parsedTime, err := time.Parse(time.RFC3339, scheduledStr); err == nil {
							data[dbField] = parsedTime
						} else {
							data[dbField] = value
						}
					} else {
						data[dbField] = scheduledStr
					}
				} else {
					data[dbField] = nil
				}
			} else {
				data[dbField] = value
			}
		}
	}
	return data
}

func buildUpdateQuery(id string, data map[string]interface{}) (string, []interface{}, error) {
	if len(data) == 0 {
		return "", nil, fmt.Errorf("no data to update")
	}

	setClauses := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data)+1)
	i := 1

	for column, value := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE drafts SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), i)

	return query, args, nil
}

func validatePublishRequest(req models.DraftDataPost) error {
	if req.Title == "" || req.Article == "" || len(req.TargetProducts) == 0 {
		return fmt.Errorf("title, article, and target_products are required")
	}
	return nil
}
