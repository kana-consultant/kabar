package history

type HistoryFilter struct {
	Status string
	Action string
	Search string
	Limit  int
	Offset int
}

type HistoryStats struct {
	Total          int     `json:"total"`
	SuccessCount   int     `json:"success"`
	FailedCount    int     `json:"failed"`
	PublishedCount int     `json:"published"`
	ScheduledCount int     `json:"scheduled"`
	SuccessRate    float64 `json:"successRate"`
}

type CreateHistoryRequest struct {
	Title          string
	Topic          string
	Content        string
	ImageURL       *string
	TargetProducts []string
	Status         string
	Action         string
	ErrorMessage   *string
	ScheduledFor   *string
	CreatedBy      string
	TeamID         string
}
