// services/sitemap_service.go
package sitemap

import (
	"context"
	"fmt"
	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/sitemap"
	"strings"
	"time"
)

type sitemapService struct {
	historyRepo history.HistoryRepository
}

// NewSitemapService creates a new sitemap service instance
func NewSitemapService(historyRepo history.HistoryRepository) sitemap.Service {
	return &sitemapService{
		historyRepo: historyRepo,
	}
}

// GenerateSitemap generates sitemap XML from published histories
func (s *sitemapService) GenerateSitemap(
	ctx context.Context,
	req sitemap.GenerateRequest,
) (*sitemap.GenerateResponse, error) {

	// 1. Get all published histories from repository
	filter := history.HistoryFilter{
		Limit:  req.Limit,
		Offset: 0,
	}

	paginatedResult, err := s.historyRepo.GetAllPublished(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get published histories: %w", err)
	}

	// 2. Convert to SitemapArticle
	articles := make([]sitemap.SitemapArticle, 0, len(paginatedResult.Data))
	for _, h := range paginatedResult.Data {
		// Generate slug from title
		slug := generateSlug(h.Title)

		// Combine baseURL + slug
		finalURL := strings.TrimSuffix(req.BaseURL, "/") + "/" + slug

		// Get image URL (if exists)
		var imageURL string
		if h.ImageURL != nil {
			imageURL = *h.ImageURL
		}

		articles = append(articles, sitemap.SitemapArticle{
			URL:       finalURL,
			Title:     h.Title,
			ImageURL:  imageURL,
			UpdatedAt: h.CreatedAt, // atau h.PublishedAt jika ada
		})
	}

	// 3. Build sitemap XML
	sitemapXML := s.buildSitemapXML(articles, req.IncludeImages)

	return &sitemap.GenerateResponse{
		SitemapXML:  sitemapXML,
		TotalURLs:   len(articles),
		GeneratedAt: time.Now(),
	}, nil
}

// buildSitemapXML builds XML sitemap from articles
func (s *sitemapService) buildSitemapXML(articles []sitemap.SitemapArticle, includeImages bool) string {
	var xml strings.Builder

	// XML header
	xml.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	// URLSet with namespace
	xml.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`)

	// Check if any article has image
	hasImage := false
	if includeImages {
		for _, a := range articles {
			if a.ImageURL != "" {
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
	for _, article := range articles {
		if article.URL == "" {
			continue
		}

		xml.WriteString(`  <url>` + "\n")
		xml.WriteString(`    <loc>` + escapeXML(article.URL) + `</loc>` + "\n")

		// Last modification date
		if !article.UpdatedAt.IsZero() {
			xml.WriteString(`    <lastmod>` + article.UpdatedAt.Format("2006-01-02") + `</lastmod>` + "\n")
		}

		// Change frequency (default weekly)
		xml.WriteString(`    <changefreq>weekly</changefreq>` + "\n")

		// Priority (default 0.8)
		xml.WriteString(`    <priority>0.8</priority>` + "\n")

		// Image (if available)
		if includeImages && article.ImageURL != "" {
			xml.WriteString(`    <image:image>` + "\n")
			xml.WriteString(`      <image:loc>` + escapeXML(article.ImageURL) + `</image:loc>` + "\n")
			if article.Title != "" {
				xml.WriteString(`      <image:title>` + escapeXML(article.Title) + `</image:title>` + "\n")
			}
			xml.WriteString(`    </image:image>` + "\n")
		}

		xml.WriteString(`  </url>` + "\n")
	}

	xml.WriteString(`</urlset>`)
	return xml.String()
}

// generateSlug generates a URL slug from title
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, ch := range slug {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			result.WriteRune(ch)
		}
	}

	// Trim leading/trailing hyphens
	return strings.Trim(result.String(), "-")
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
