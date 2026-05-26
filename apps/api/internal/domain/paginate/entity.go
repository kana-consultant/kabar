package paginate

type PaginationParams struct {
	Limit  int
	Offset int
}

type PaginatedResult[T any] struct {
	Data        []T `json:"data"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`

	TotalSuccess int `json:"total_success,omitempty"`
	TotalFailed  int `json:"total_failed,omitempty"`
}
