package helper

import (
	"seo-backend/internal/domain/draft"
	"strconv"
	"strings"
	"time"
)

func generateSitemapInfo(sitemapConfig map[string]interface{}, draft draft.DraftDataPost, baseURL string) map[string]interface{} {
	if sitemapConfig == nil {
		return nil
	}

	if enabled, ok := sitemapConfig["enabled"].(bool); !ok || !enabled {
		return nil
	}

	dynamicConfig, ok := sitemapConfig["dynamicConfig"].(map[string]interface{})
	if !ok {
		return nil
	}

	// URL
	urlPattern := getStringValue(dynamicConfig, "urlPattern", "/p/{id}")
	urlPath := replaceAllPlaceholders(urlPattern, draft)
	fullURL := baseURL + urlPath

	// Priority → float64, clamp 0.0–1.0
	prioritySource := getStringValue(dynamicConfig, "prioritySource", "0.7")
	priorityStr := getValueFromDraftWithPlaceholder(draft, prioritySource)
	if priorityStr == "" {
		priorityStr = "0.7"
	}
	priorityFloat, err := strconv.ParseFloat(priorityStr, 64)
	if err != nil {
		priorityFloat = 0.7
	}
	if priorityFloat < 0.0 {
		priorityFloat = 0.0
	}
	if priorityFloat > 1.0 {
		priorityFloat = 1.0
	}

	// Changefreq
	changefreqSource := getStringValue(dynamicConfig, "changefreqSource", "weekly")
	changefreq := getValueFromDraftWithPlaceholder(draft, changefreqSource)
	if changefreq == "" {
		changefreq = "weekly"
	}

	// Lastmod → ISO 8601
	lastmodSource := getStringValue(dynamicConfig, "lastmodSource", "updated_at")
	lastmod := getValueFromDraftWithPlaceholder(draft, lastmodSource)
	if lastmod == "" {
		lastmod = time.Now().Format("2006-01-02") // YYYY-MM-DD
	}

	// Image → harus absolute URL
	imageSource := getStringValue(dynamicConfig, "imageSource", "image_url")
	imageURL := getValueFromDraftWithPlaceholder(draft, imageSource)

	sitemapInfo := map[string]interface{}{
		"loc":        fullURL,
		"lastmod":    lastmod,
		"changefreq": changefreq,
		"priority":   priorityFloat,
	}

	if imageURL != "" {
		if !strings.HasPrefix(imageURL, "http") {
			imageURL = baseURL + imageURL
		}
		sitemapInfo["image"] = map[string]string{
			"loc": imageURL,
		}
	}

	// static_urls TIDAK dimasukkan di sini
	// handle terpisah di caller

	return sitemapInfo
}
