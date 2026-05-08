package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
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
			ParseWIBTime(scheduledFor)

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

func ParseWIBTime(timeStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}, err
	}

	// 1. Kalau ada timezone (Z / +07:00), parse langsung
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t.In(loc), nil
	}

	// 2. Format TANPA timezone → anggap WIB
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, loc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
}

func WriteAllProductsFailed(
	w http.ResponseWriter,
	result interface{},
) {

	response := map[string]interface{}{
		"message": "Failed to publish draft: all products failed",
		"status":  "failed",
		"results": result,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusBadGateway)

	json.NewEncoder(w).Encode(response)
}
