package helper

import (
	"seo-backend/internal/domain/draft"
	"time"
)

func generateSitemapInfo(sitemapConfig map[string]interface{}, draft draft.DraftDataPost, baseURL string) map[string]interface{} {
	if sitemapConfig == nil {
		return nil
	}

	// Check enabled
	if enabled, ok := sitemapConfig["enabled"].(bool); !ok || !enabled {
		return nil
	}

	// Get dynamic config
	dynamicConfig, ok := sitemapConfig["dynamicConfig"].(map[string]interface{})
	if !ok {
		return nil
	}

	// URL Pattern
	urlPattern := getStringValue(dynamicConfig, "urlPattern", "/p/{id}")
	urlPath := replaceAllPlaceholders(urlPattern, draft)
	fullURL := baseURL + urlPath

	// Priority
	prioritySource := getStringValue(dynamicConfig, "prioritySource", "0.7")
	priority := getValueFromDraftWithPlaceholder(draft, prioritySource)
	if priority == "" {
		priority = "0.7"
	}

	// Changefreq
	changefreqSource := getStringValue(dynamicConfig, "changefreqSource", "weekly")
	changefreq := getValueFromDraftWithPlaceholder(draft, changefreqSource)
	if changefreq == "" {
		changefreq = "weekly"
	}

	// Image
	imageSource := getStringValue(dynamicConfig, "imageSource", "image_url")
	imageURL := getValueFromDraftWithPlaceholder(draft, imageSource)

	// Lastmod
	lastmodSource := getStringValue(dynamicConfig, "lastmodSource", "updated_at")
	lastmod := getValueFromDraftWithPlaceholder(draft, lastmodSource)
	if lastmod == "" {
		lastmod = time.Now().Format(time.RFC3339)
	}

	sitemapInfo := map[string]interface{}{
		"loc":        fullURL,
		"lastmod":    lastmod,
		"changefreq": changefreq,
		"priority":   priority,
	}

	if imageURL != "" {
		sitemapInfo["image"] = map[string]string{
			"loc": imageURL,
		}
	}

	// Static URLs
	if staticUrls, ok := sitemapConfig["staticUrls"].([]interface{}); ok {
		sitemapInfo["static_urls"] = staticUrls
	}

	return sitemapInfo
}
