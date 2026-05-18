package helper

import (
	"seo-backend/internal/domain/draft"
	"strings"
	"time"
)

func replaceAllPlaceholders(text string, draft draft.DraftDataPost) string {
	plainContent := stripHTML(draft.Article)

	// Basic placeholders
	text = strings.ReplaceAll(text, "{title}", draft.Title)
	text = strings.ReplaceAll(text, "{topic}", draft.Topic)
	text = strings.ReplaceAll(text, "{content}", draft.Article)

	// Excerpt
	excerpt := plainContent
	if len(excerpt) > 160 {
		excerpt = excerpt[:160] + "..."
	}
	text = strings.ReplaceAll(text, "{excerpt}", excerpt)

	// Image URL
	if draft.ImageURL != nil {
		text = strings.ReplaceAll(text, "{image_url}", *draft.ImageURL)
	} else {
		text = strings.ReplaceAll(text, "{image_url}", "")
	}

	// Meta placeholders
	text = strings.ReplaceAll(text, "{meta_title}", draft.Title)
	text = strings.ReplaceAll(text, "{meta_description}", excerpt)
	text = strings.ReplaceAll(text, "{meta_keywords}", draft.Topic)

	// OG placeholders
	text = strings.ReplaceAll(text, "{og_title}", draft.Title)
	text = strings.ReplaceAll(text, "{og_description}", excerpt)
	if draft.ImageURL != nil {
		text = strings.ReplaceAll(text, "{og_image}", *draft.ImageURL)
	} else {
		text = strings.ReplaceAll(text, "{og_image}", "")
	}

	// Twitter placeholders
	text = strings.ReplaceAll(text, "{twitter_title}", draft.Title)
	text = strings.ReplaceAll(text, "{twitter_description}", excerpt)
	if draft.ImageURL != nil {
		text = strings.ReplaceAll(text, "{twitter_image}", *draft.ImageURL)
	} else {
		text = strings.ReplaceAll(text, "{twitter_image}", "")
	}

	// Sitemap placeholders
	text = strings.ReplaceAll(text, "{sitemap_priority}", "0.7")
	text = strings.ReplaceAll(text, "{sitemap_changefreq}", "weekly")

	// Timestamp
	text = strings.ReplaceAll(text, "{timestamp}", time.Now().Format(time.RFC3339))

	return text
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
