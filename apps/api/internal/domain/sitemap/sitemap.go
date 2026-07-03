// domain/sitemap/service.go
package sitemap

import (
	"context"
	"time"
)

// Service defines the sitemap service interface
type Service interface {
	// GenerateSitemap generates sitemap XML from published histories
	GenerateSitemap(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
}

// GenerateRequest represents request to generate sitemap
type GenerateRequest struct {
	BaseURL       string // Base URL of the website (e.g., https://client.com)
	IncludeImages bool   // Whether to include images in sitemap
	Limit         int    // Max number of articles (0 = all)
}

// GenerateResponse represents response from sitemap generation
type GenerateResponse struct {
	SitemapXML  string    // Generated sitemap XML content
	TotalURLs   int       // Total number of URLs included
	GeneratedAt time.Time // Generation timestamp
}

// SitemapArticle represents article data needed for sitemap
type SitemapArticle struct {
	URL       string    // Final URL of the article
	Title     string    // Article title
	ImageURL  string    // Image URL (optional)
	UpdatedAt time.Time // Last update time
}
