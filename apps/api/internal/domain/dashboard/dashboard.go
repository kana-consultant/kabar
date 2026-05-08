// DashboardStats represents the dashboard statistics
package dashboard

type DashboardStats struct {
	TotalContent        int     `json:"totalContent"`
	TotalProducts       int     `json:"totalProducts"`
	TotalPublished      int     `json:"totalPublished"`
	AverageSeoScore     float64 `json:"averageSeoScore"`
	ContentChange       string  `json:"contentChange"`
	ProductsChange      string  `json:"productsChange"`
	PublishedPercentage int     `json:"publishedPercentage"`
	SeoScoreChange      string  `json:"seoScoreChange"`
}

// UserContext for dashboard
type DashboardUserContext struct {
	UserID string
	TeamID string
	Role   string
}
