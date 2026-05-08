package models

type GenerateArticleRequest struct {
	Topic    string   `json:"topic"`
	Tone     string   `json:"tone,omitempty"`
	Length   string   `json:"length,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
	Language string   `json:"language,omitempty"`
}

type GenerateArticleResponse struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Excerpt     string   `json:"excerpt"`
	Keywords    []string `json:"keywords"`
	ImagePrompt string   `json:"imagePrompt"`
	ImageURL    string   `json:"imageUrl"`
}

type GenerateImageRequest struct {
	Prompt string `json:"prompt"`
}

type GenerateImageResponse struct {
	ImageURL string `json:"imageUrl"`
}