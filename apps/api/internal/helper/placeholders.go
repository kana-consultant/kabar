package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"seo-backend/internal/domain/draft"
	"strings"
	"time"

	"github.com/google/uuid"
)

func stripHTMLTags(html string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return strings.TrimSpace(re.ReplaceAllString(html, ""))
}

func injectImageAfterH1(content string, imageTag string) string {
	if imageTag == "" {
		return content
	}
	idx := strings.Index(content, "</h1>")
	if idx == -1 {
		// tidak ada <h1>, taruh di awal
		return imageTag + content
	}
	insertAt := idx + len("</h1>")
	return content[:insertAt] + imageTag + content[insertAt:]
}

func replaceAllPlaceholders(text string, draft draft.DraftDataPost) string {

	// Excerpt
	excerpt := draft.Excerpt

	// Image URL
	imageURL := ""
	imageTag := ""
	if draft.ImageURL != nil {
		imageURL = *draft.ImageURL
		imageTag = fmt.Sprintf(`<img src="%s" wrapperstyle="display: flex">`, imageURL)
	}

	// Use existing ID or generate new UUID
	generatedID := draft.Id
	if generatedID == "" {
		generatedID = uuid.New().String()
	}

	keywordNames := make([]string, len(draft.Keywords))
	copy(keywordNames, draft.Keywords)

	keywordsJSON, _ := json.Marshal(keywordNames)

	slug := draft.Slug
	if slug == "" {
		base := draft.Title
		if base == "" {
			base = draft.Topic
		}
		slug = slugify(base)
	}

	// Content variants
	contentHTML := draft.Article
	contentText := stripHTMLTags(draft.Article)
	contentWithImage := injectImageAfterH1(draft.Article, imageTag)

	placeholders := map[string]string{
		"{title}":              draft.Title,
		"{topic}":              draft.Topic,
		"{content}":            contentHTML,      // HTML asli
		"{content_text}":       contentText,      // plain text (tanpa tag HTML)
		"{content_with_image}": contentWithImage, // HTML + gambar setelah <h1>
		"{slug}":               slug,
		"{excerpt}":            excerpt,
		"{image_url}":          imageURL, // plain URL
		"{image_content_html}": imageTag, // <img> tag saja

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
func getArrayPlaceholders(draft draft.DraftDataPost) map[string][]string {
	keywordNames := make([]string, len(draft.Keywords))
	copy(keywordNames, draft.Keywords)

	return map[string][]string{
		"{tags}":     keywordNames,
		"{keywords}": keywordNames,
	}
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
		return stripHTML(draft.Title)
	case "topic":
		return stripHTML(draft.Topic)
	case "content", "article":
		return draft.Article
	case "excerpt":
		excerpt := stripHTML(draft.Article)
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
		return stripHTML(draft.Title)
	case "meta_description", "og_description", "twitter_description":
		excerpt := stripHTML(draft.Article)
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
