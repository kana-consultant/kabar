package helper

import (
	"seo-backend/internal/domain/draft"
	"time"
)

func generateMetaTags(metaConfig map[string]interface{}, draft draft.DraftDataPost, baseURL string, sitemapConfig map[string]interface{}) map[string]string {
	result := make(map[string]string)

	if metaConfig == nil {
		return result
	}

	// Check enabled
	if enabled, ok := metaConfig["enabled"].(bool); !ok || !enabled {
		return result
	}

	// Get default tags
	if defaultTags, ok := metaConfig["defaultTags"].(map[string]interface{}); ok {
		if charset, ok := defaultTags["charset"].(string); ok {
			result["charset"] = charset
		}
		if viewport, ok := defaultTags["viewport"].(string); ok {
			result["viewport"] = viewport
		}
		if robots, ok := defaultTags["robots"].(string); ok {
			result["robots"] = robots
		}
		if generator, ok := defaultTags["generator"].(string); ok {
			result["generator"] = generator
		}
	}

	// Get dynamic tags
	dynamicTags, ok := metaConfig["dynamicTags"].(map[string]interface{})
	if !ok {
		return result
	}

	// Get sources
	titleSource := getStringValue(dynamicTags, "titleSource", "{title}")
	descSource := getStringValue(dynamicTags, "descriptionSource", "{excerpt}")
	imageSource := getStringValue(dynamicTags, "imageSource", "{image_url}")

	// Get values
	title := getValueFromDraftWithPlaceholder(draft, titleSource)
	description := getValueFromDraftWithPlaceholder(draft, descSource)
	imageURL := getValueFromDraftWithPlaceholder(draft, imageSource)

	// Basic meta tags
	if title != "" {
		result["title"] = title
		result["og:title"] = title
		result["twitter:title"] = title
	}

	if description != "" {
		result["description"] = description
		result["og:description"] = description
		result["twitter:description"] = description
	}

	if imageURL != "" {
		result["og:image"] = imageURL
		result["twitter:image"] = imageURL
	}

	// Standard social media tags
	result["og:type"] = "article"
	result["og:site_name"] = getStringValue(dynamicTags, "siteName", "AI Content Generator")
	result["twitter:card"] = "summary_large_image"

	// LinkedIn tags
	if title != "" {
		result["linkedin:title"] = title
	}
	if description != "" {
		result["linkedin:description"] = description
	}
	if imageURL != "" {
		result["linkedin:image"] = imageURL
	}

	// Pinterest tags
	if imageURL != "" {
		result["pinterest:image"] = imageURL
	}
	if description != "" {
		result["pinterest:description"] = description
	}

	// Canonical URL
	canonicalURL := getCanonicalURL(draft, baseURL, sitemapConfig)
	if canonicalURL != "" {
		result["canonical"] = canonicalURL
		result["og:url"] = canonicalURL
	}

	// Custom tags
	if customTags, ok := dynamicTags["customTags"].(map[string]interface{}); ok {
		for key, value := range customTags {
			if strValue, ok := value.(string); ok {
				if replacedValue := replaceAllPlaceholders(strValue, draft); replacedValue != "" {
					result[key] = replacedValue
				}
			}
		}
	}

	// Cleanup empty values
	for key, value := range result {
		if value == "" {
			delete(result, key)
		}
	}

	return result
}

func getCanonicalURL(draft draft.DraftDataPost, baseURL string, sitemapConfig map[string]interface{}) string {
	// Try from sitemap config
	if sitemapConfig != nil {
		if dynamicConfig, ok := sitemapConfig["dynamicConfig"].(map[string]interface{}); ok {
			if urlPattern, ok := dynamicConfig["urlPattern"].(string); ok {
				if canonical := replaceAllPlaceholders(urlPattern, draft); canonical != "" {
					return baseURL + canonical
				}
			}
		}
	}

	// Default fallback
	return baseURL + "/p/" + time.Now().Format("20060102150405")
}
