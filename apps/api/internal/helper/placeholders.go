package helper

import (
	"encoding/json"
	"seo-backend/internal/domain/draft"
	"strings"
	"time"

	"github.com/google/uuid"
)

func replaceAllPlaceholders(text string, draft draft.DraftDataPost) string {
	plainContent := stripHTML(draft.Article)

	// Excerpt
	excerpt := plainContent
	if len(excerpt) > 160 {
		excerpt = excerpt[:160] + "..."
	}

	// Image URL
	imageURL := ""
	if draft.ImageURL != nil {
		imageURL = *draft.ImageURL
	}

	// Use existing ID or generate new UUID
	generatedID := draft.Id
	if generatedID == "" {
		generatedID = uuid.New().String()
	}

	keywordNames := make([]string, len(draft.Keywords))

	keywordsJSON, _ := json.Marshal(keywordNames)
	keywordsStr := strings.Join(keywordNames, ", ") // ✅ convert ke string dulu

	slug := draft.Slug
	if slug == "" {
		base := draft.Title
		if base == "" {
			base = draft.Topic
		}
		slug = slugify(base)
	}

	placeholders := map[string]string{
		"{title}":     draft.Title,
		"{topic}":     draft.Topic,
		"{content}":   draft.Article,
		"{slug}":      slug,
		"{tags}":      keywordsStr, // ✅ sudah string
		"{excerpt}":   excerpt,
		"{image_url}": imageURL,

		// Meta
		"{meta_title}":       draft.Title,
		"{meta_description}": excerpt,
		"{meta_keywords}":    string(keywordsJSON),

		// OG
		"{og_title}":       draft.Title,
		"{og_description}": excerpt,
		"{og_image}":       imageURL,

		// Twitter
		"{twitter_title}":       draft.Title,
		"{twitter_description}": excerpt,
		"{twitter_image}":       imageURL,

		// Sitemap
		"{sitemap_priority}":   "0.7",
		"{sitemap_changefreq}": "weekly",

		// Timestamp
		"{timestamp}": time.Now().Format(time.RFC3339),

		// ID
		"{id}": generatedID,
	}

	for key, value := range placeholders {
		text = strings.ReplaceAll(text, key, value)
	}

	return text
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")

	// Hapus karakter selain huruf, angka, dan tanda hubung
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}

	// Hilangkan tanda hubung ganda
	result := b.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}

	return strings.Trim(result, "-")
}

func getValueFromDraftWithPlaceholder(draft draft.DraftDataPost, source string) string {
	if strings.HasPrefix(source, "{") && strings.HasSuffix(source, "}") {
		return replaceAllPlaceholders(source, draft)
	}

	switch source {
	case "title":
		return draft.Title
	case "topic":
		return draft.Topic
	case "content", "article":
		return draft.Article
	case "excerpt":
		excerpt := draft.Article
		if len(excerpt) > 160 {
			excerpt = excerpt[:160] + "..."
		}
		return excerpt
	case "image_url":
		if draft.ImageURL != nil {
			return *draft.ImageURL
		}
		return ""
	case "meta_title", "og_title", "twitter_title":
		return draft.Title
	case "meta_description", "og_description", "twitter_description":
		excerpt := draft.Article
		if len(excerpt) > 160 {
			excerpt = excerpt[:160] + "..."
		}
		return excerpt
	case "og_image", "twitter_image":
		if draft.ImageURL != nil {
			return *draft.ImageURL
		}
		return ""
	case "sitemap_priority":
		return "0.7"
	case "sitemap_changefreq":
		return "weekly"
	case "updated_at", "created_at":
		return time.Now().Format(time.RFC3339)
	default:
		return ""
	}
}
