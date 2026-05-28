package userRole

import (
	"fmt"
	"seo-backend/internal/models"
	"strings"
)

// BuildAccessFilter builds WHERE clause based on user role for data access control
func BuildAccessFilter(filter models.UserContext) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	switch filter.GetRole() {
	case "superadmin":
		conditions = append(conditions, "1=1")

	case "admin":
		if filter.GetTeamID() != "" && filter.GetTeamID() != "00000000-0000-0000-0000-000000000000" {
			conditions = append(conditions, fmt.Sprintf("team_id = $%d", argIndex))
			args = append(args, filter.GetTeamID())
			argIndex++
		} else {
			conditions = append(conditions, fmt.Sprintf("created_by = $%d", argIndex))
			args = append(args, filter.GetUserID())
			argIndex++
		}

	default:
		conditions = append(conditions, fmt.Sprintf("team_id = $%d", argIndex))
		args = append(args, filter.GetTeamID())
	}

	if len(conditions) == 0 {
		return "1=1", nil
	}

	return strings.Join(conditions, " AND "), args
}
