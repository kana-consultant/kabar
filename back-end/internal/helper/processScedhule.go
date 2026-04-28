package helper

import (
	"fmt"
	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"strings"
	"time"
)

func ScheduleDraft(
	id string,
	scheduledFor string,
) bool {

	data := map[string]interface{}{
		"status":     "scheduled",
		"updated_at": time.Now(),
	}

	if strings.HasPrefix(
		scheduledFor,
		"daily:",
	) {

		data["scheduled_pattern"] = scheduledFor

	} else {

		parsedTime, err :=
			parseTimeFlexible(scheduledFor)

		if err != nil {
			return false
		}

		data["scheduled_for"] = parsedTime
	}

	qb := builder.NewQueryBuilder("drafts")

	sqlQuery, args, err :=
		qb.Update().
			SetMap(data).
			WhereEq("id", id).
			Build()

	if err != nil {
		return false
	}

	_, err = database.GetDB().
		Exec(sqlQuery, args...)

	if err != nil {
		return false
	}

	return true
}

// Helper function to parse time flexibly
func parseTimeFlexible(timeStr string) (time.Time, error) {
	// List of possible formats
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // Without timezone
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, timeStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}
