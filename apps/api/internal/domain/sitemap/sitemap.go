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

	// GetSitemapHistory retrieves sitemap generation history
	GetSitemapHistory(ctx context.Context) ([]HistoryResponse, error)
}

// GenerateRequest represents request to generate sitemap
type GenerateRequest struct {
	ProductID     string // ID produk yang dipilih
	BaseURL       string // Base URL of the website (e.g., https://client.com)
	IncludeImages bool   // Whether to include images in sitemap
}

// GenerateResponse represents response from sitemap generation
type GenerateResponse struct {
	SitemapXML  string    // Generated sitemap XML content
	TotalURLs   int       // Total number of URLs included
	GeneratedAt time.Time // Generation timestamp
	ProductID   string    // Product ID used
	BaseURL     string    // Base URL used
}

// HistoryResponse represents sitemap history item
type HistoryResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	TotalURLs   int       `json:"totalURLs"`
	Status      string    `json:"status"` // success, failed, pending
	CreatedAt   time.Time `json:"createdAt"`
	SitemapURL  string    `json:"sitemapURL,omitempty"`
	ProductID   string    `json:"productId,omitempty"`
	ProductName string    `json:"productName,omitempty"`
	BaseURL     string    `json:"baseURL,omitempty"`
}
