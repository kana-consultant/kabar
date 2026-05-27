// internal/domain/generate/entity.go
package generate

import "time"

// ========== PARAMETERS ==========
type ArticleGenerationParams struct {
	Topic             string
	ModelID           string
	Tone              string
	Length            string
	Language          string
	AutoGenerateImage bool
}

type ImageGenerationParams struct {
	Prompt  string
	ModelID string
}

// ========== RESULTS ==========
type ArticleResult struct {
	Title            string   `json:"title"`
	Content          string   `json:"content"`
	Excerpt          string   `json:"excerpt"`
	Keywords         []string `json:"keywords"`
	ImagePrompt      string   `json:"imagePrompt"`
	ImageURL         string   `json:"imageUrl"`
	WordCount        int      `json:"wordCount"`
	ReadabilityScore int      `json:"readabilityScore"`
	SeoScore         int      `json:"seo_score"`
	Slug             string   `json:"slug"`
}

type ImageResult struct {
	ImageURL    string `json:"imageUrl"`
	Prompt      string `json:"prompt"`
	GeneratedAt string `json:"generatedAt"`
	Model       string `json:"model"`
}

// ========== INTERNAL TYPES ==========
type ModelConfig struct {
	APIKey            string
	ModelName         string
	Template          string
	ResponsePath      string
	ResponseImagePath string
	BaseURL           string
	AuthType          string
	AuthHeader        string
	AuthPrefix        string
	Endpoint          string
	SystemPrompt      string
	EncryptedKey      string
	Slug              string
}

type GenerationHistory struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // "article" or "image"
	Topic     string    `json:"topic"`
	Prompt    string    `json:"prompt"`
	Result    string    `json:"result"`
	ModelID   string    `json:"model_id"`
	CreatedAt time.Time `json:"created_at"`
}
