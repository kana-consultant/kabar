package helper

import (
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
