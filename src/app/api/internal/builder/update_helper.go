package builder

import (
	"fmt"
	"sort"
	"strings"
)

// buildUpdateByIDQuery builds dynamic UPDATE query by ID
func BuildUpdateByIDQuery(table string, id string, data map[string]interface{}) (string, []interface{}, error) {
	if len(data) == 0 {
		return "", nil, fmt.Errorf("empty update data")
	}

	setParts := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data)+1)

	i := 1

	// supaya deterministik (optional tapi bagus)
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, col := range keys {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, data[col])
		i++
	}

	// WHERE id
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		table,
		strings.Join(setParts, ", "),
		i,
	)

	args = append(args, id)

	return query, args, nil
}
