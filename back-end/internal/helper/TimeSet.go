// helper/time.go
package helper

import "time"

func ParseWIBTime(dateStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation(
		"2006-01-02T15:04:05",
		dateStr,
		loc,
	)
}
