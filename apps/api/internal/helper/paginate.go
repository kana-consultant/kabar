package helper

import (
	"net/http"
	"seo-backend/internal/domain/paginate"
	"strconv"
)

func ParsePaginationParams(r *http.Request) paginate.PaginationParams {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	return paginate.PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}
