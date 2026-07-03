package helper

import (
	"encoding/json"
	"net/http"
	"time"
)

func ParseWIBTime(timeStr string) time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}
	}

	// 1. Kalau ada timezone (Z / +07:00), parse langsung
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t.In(loc)
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
			return t
		}
	}

	return time.Time{}
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
