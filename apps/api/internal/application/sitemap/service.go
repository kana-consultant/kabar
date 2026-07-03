// internal/application/sitemap/service.go
package sitemap

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/domain/sitemap"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	historyRepo history.HistoryRepository
	productRepo product.ProductRepository
}

func NewService(
	historyRepo history.HistoryRepository,
	productRepo product.ProductRepository,
) sitemap.Service {
	return &Service{
		historyRepo: historyRepo,
		productRepo: productRepo,
	}
}

// GenerateSitemap generates sitemap XML from published histories
func (s *Service) GenerateSitemap(
	ctx context.Context,
	req sitemap.GenerateRequest,
) (*sitemap.GenerateResponse, error) {

	log.Printf("[GenerateSitemap] START - productID=%s, baseURL=%s", req.ProductID, req.BaseURL)

	// 1. Get product data by ID
	productData, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		log.Printf("[GenerateSitemap] ERROR failed to get product: %v", err)
		return nil, fmt.Errorf("failed to get product %s: %w", req.ProductID, err)
	}

	log.Printf("[GenerateSitemap] Product found: ID=%s, Name=%s", productData.ID, productData.Name)

	// 2. Get all published histories for this product
	filter := history.HistoryFilter{
		Limit:     req.Limit,
		Offset:    0,
		ProductID: req.ProductID,
	}

	paginatedResult, err := s.historyRepo.GetAllPublished(ctx, filter)
	if err != nil {
		log.Printf("[GenerateSitemap] ERROR failed to get published histories: %v", err)
		return nil, fmt.Errorf("failed to get published histories for product %s: %w", req.ProductID, err)
	}

	log.Printf("[GenerateSitemap] Got %d histories for product %s", len(paginatedResult.Data), req.ProductID)

	// 3. Build sitemap XML from histories with placeholder support
	sitemapXML := s.buildSitemapXML(paginatedResult.Data, req, productData)

	response := &sitemap.GenerateResponse{
		SitemapXML:  sitemapXML,
		TotalURLs:   len(paginatedResult.Data),
		GeneratedAt: time.Now(),
		ProductID:   req.ProductID,
		BaseURL:     req.BaseURL,
	}

	log.Printf("[GenerateSitemap] SUCCESS - generated sitemap with %d URLs", response.TotalURLs)

	return response, nil
}

// GetSitemapHistory retrieves sitemap generation history
func (s *Service) GetSitemapHistory(ctx context.Context) ([]sitemap.HistoryResponse, error) {
	log.Printf("[GetSitemapHistory] START fetching history")

	filter := history.HistoryFilter{
		Limit:  50,
		Offset: 0,
	}

	log.Printf("[GetSitemapHistory] Filter: limit=%d, offset=%d", filter.Limit, filter.Offset)

	paginatedResult, err := s.historyRepo.GetAllPublished(ctx, filter)
	if err != nil {
		log.Printf("[GetSitemapHistory] ERROR failed to get published histories: %v", err)
		return nil, fmt.Errorf("failed to get sitemap history: %w", err)
	}

	log.Printf("[GetSitemapHistory] SUCCESS got %d histories (total: %d)",
		len(paginatedResult.Data), paginatedResult.TotalItems)

	// Log sample data
	if len(paginatedResult.Data) > 0 {
		log.Printf("[GetSitemapHistory] Sample data - first history: ID=%s, Title=%s, Status=%s, CreatedAt=%v",
			paginatedResult.Data[0].ID,
			paginatedResult.Data[0].Title,
			paginatedResult.Data[0].Status,
			paginatedResult.Data[0].CreatedAt,
		)

		// Log target_products for first item
		if len(paginatedResult.Data[0].TargetProducts) > 0 {
			log.Printf("[GetSitemapHistory] Sample target_products: %v",
				paginatedResult.Data[0].TargetProducts)
		}
	}

	responses := make([]sitemap.HistoryResponse, 0, len(paginatedResult.Data))

	log.Printf("[GetSitemapHistory] Processing %d histories to response", len(paginatedResult.Data))

	for i, h := range paginatedResult.Data {
		status := "success"
		if h.Status == "failed" {
			status = "failed"
		}

		// Get product name from first target product
		productID := ""
		productName := ""
		baseURL := ""

		if len(h.TargetProducts) > 0 {
			productID = h.TargetProducts[0]

			// Get product data to get name and baseURL
			productData, err := s.productRepo.GetByID(ctx, productID)
			if err == nil && productData != nil {
				productName = productData.Name
				// Use api_endpoint or domain from product
				if productData.APIEndpoint != "" {
					baseURL = productData.APIEndpoint
				}
				log.Printf("[GetSitemapHistory] Product data: ID=%s, Name=%s, BaseURL=%s",
					productID, productName, baseURL)
			} else {
				log.Printf("[GetSitemapHistory] WARNING failed to get product %s: %v", productID, err)
			}
		}

		// Generate sitemap URL
		sitemapURL := ""
		if baseURL != "" && productID != "" {
			sitemapURL = fmt.Sprintf("/sitemap?product_id=%s&base_url=%s", productID, baseURL)
		}

		responses = append(responses, sitemap.HistoryResponse{
			ID:          h.ID,
			Title:       h.Title,
			TotalURLs:   1,
			Status:      status,
			CreatedAt:   h.CreatedAt,
			SitemapURL:  sitemapURL,
			ProductID:   productID,
			ProductName: productName,
			BaseURL:     baseURL,
		})

		// Log progress untuk first 5 items
		if i < 5 {
			log.Printf("[GetSitemapHistory] Processing item %d: ID=%s, Title=%s, ProductID=%s, BaseURL=%s",
				i, h.ID, h.Title, productID, baseURL)
		}
	}

	log.Printf("[GetSitemapHistory] SUCCESS processed %d histories into response", len(responses))

	// Log final summary
	log.Printf("[GetSitemapHistory] SUMMARY: total_histories=%d, response_count=%d, statuses=%v",
		len(paginatedResult.Data),
		len(responses),
		getStatusSummary(responses))

	return responses, nil
}

// Helper untuk summary status
func getStatusSummary(responses []sitemap.HistoryResponse) map[string]int {
	summary := make(map[string]int)
	for _, r := range responses {
		summary[r.Status]++
	}
	return summary
}

// buildSitemapXML builds XML sitemap from histories with placeholder support
func (s *Service) buildSitemapXML(histories []history.History, req sitemap.GenerateRequest, productData *product.Product) string {
	var xml strings.Builder

	// XML header
	xml.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	// URLSet with namespace
	xml.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`)

	// Check if any history has image
	hasImage := false
	if req.IncludeImages {
		for _, h := range histories {
			if h.ImageURL != nil && *h.ImageURL != "" {
				hasImage = true
				break
			}
		}
	}

	if hasImage {
		xml.WriteString(` xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"`)
	}

	xml.WriteString(`>` + "\n")

	// Add each URL
	baseURL := strings.TrimSuffix(req.BaseURL, "/")

	for _, h := range histories {
		// Generate URL with placeholder replacement
		finalURL := s.replacePlaceholders(baseURL, h, productData)

		xml.WriteString(`  <url>` + "\n")
		xml.WriteString(`    <loc>` + escapeXML(finalURL) + `</loc>` + "\n")

		// Last modification date
		if !h.CreatedAt.IsZero() {
			xml.WriteString(`    <lastmod>` + h.CreatedAt.Format("2006-01-02") + `</lastmod>` + "\n")
		}

		// Change frequency
		xml.WriteString(`    <changefreq>weekly</changefreq>` + "\n")

		// Priority
		xml.WriteString(`    <priority>0.8</priority>` + "\n")

		// Image (if available)
		if req.IncludeImages && h.ImageURL != nil && *h.ImageURL != "" {
			xml.WriteString(`    <image:image>` + "\n")
			xml.WriteString(`      <image:loc>` + escapeXML(*h.ImageURL) + `</image:loc>` + "\n")
			if h.Title != "" {
				xml.WriteString(`      <image:title>` + escapeXML(h.Title) + `</image:title>` + "\n")
			}
			xml.WriteString(`    </image:image>` + "\n")
		}

		xml.WriteString(`  </url>` + "\n")
	}

	xml.WriteString(`</urlset>`)
	return xml.String()
}

// replacePlaceholders replaces placeholders in URL template with actual data
func (s *Service) replacePlaceholders(template string, h history.History, productData *product.Product) string {
	result := template

	// Replace history-related placeholders
	result = strings.ReplaceAll(result, "{slug}", h.Slug)
	result = strings.ReplaceAll(result, "{title}", h.Title)

	// Replace date placeholders
	if !h.CreatedAt.IsZero() {
		result = strings.ReplaceAll(result, "{date}", h.CreatedAt.Format("2006-01-02"))
		result = strings.ReplaceAll(result, "{timestamp}", strconv.FormatInt(h.CreatedAt.Unix(), 10))
	}

	// Handle any remaining placeholders that couldn't be resolved
	// Remove them or replace with empty string
	re := regexp.MustCompile(`\{[^}]+\}`)
	result = re.ReplaceAllString(result, "")

	// Clean up any double slashes (but preserve protocol)
	result = s.cleanURL(result)

	return result
}

// getCategorySlug generates a slug from category name
func (s *Service) getCategorySlug(category string) string {
	if category == "" {
		return "uncategorized"
	}
	return generateSlug(category)
}

// cleanURL removes double slashes and fixes URL formatting
func (s *Service) cleanURL(url string) string {
	// Remove any double slashes except after protocol
	re := regexp.MustCompile(`([^:])//+`)
	url = re.ReplaceAllString(url, "$1/")

	// Remove trailing slash if exists
	url = strings.TrimSuffix(url, "/")

	return url
}

// generateSlug generates a URL slug from title
func generateSlug(title string) string {
	if title == "" {
		return "untitled"
	}

	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters, keep only alphanumeric and dash
	var result strings.Builder
	for _, ch := range slug {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			result.WriteRune(ch)
		}
	}

	// Remove multiple consecutive dashes
	slug = result.String()
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim dashes from start and end
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
	}

	return slug
}

// escapeXML escapes special XML characters
func escapeXML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}

// GetAvailablePlaceholders returns list of available placeholders
func (s *Service) GetAvailablePlaceholders() []string {
	return []string{
		"{slug}",
		"{title}",
		"{id}",
		"{sku}",
		"{category}",
		"{date}",
		"{timestamp}",
	}
}
