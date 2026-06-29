// internal/domain/generate/entity.go
package generate

import (
	"database/sql"
	"time"
)

// ========== PARAMETERS ==========
type ArticleGenerationParams struct {
	Topic             string `json:"topic"`
	ModelID           string `json:"modelId"`
	Tone              string `json:"tone"`
	Length            string `json:"length"`
	Language          string `json:"language"`
	Slug              string `json:"slug,omitempty"`
	ArticleID         string `json:"articleId,omitempty"`
	AutoGenerateImage bool   `json:"autoGenerateImage,omitempty"`
	ImageModelID      string `json:"imageModelId,omitempty"`
}

type ImageGenerationParams struct {
	Prompt    string `json:"prompt"`
	ModelID   string `json:"modelId"`
	Slug      string `json:"slug,omitempty"`
	ArticleID string `json:"articleId,omitempty"`
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
	APIKey            string         `json:"api_key"`
	APISystemPrompt   string         `json:"api_system_prompt"`
	ModelSystemPrompt string         `json:"model_system_prompt"`
	SystemPrompt      string         `json:"system_prompt"`
	ModelName         string         `json:"model_name"`
	MaxTokens         int            `json:"max_tokens"`
	Temperature       float64        `json:"temperature"`
	Template          string         `json:"template"`
	ResponsePath      string         `json:"response_path"`
	ResponseImagePath string         `json:"response_image_path"`
	BaseURL           string         `json:"base_url"`
	AuthType          sql.NullString `json:"auth_type"`
	AuthHeader        sql.NullString `json:"auth_header"`
	AuthPrefix        sql.NullString `json:"auth_prefix"`
	Endpoint          string         `json:"endpoint"`
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
