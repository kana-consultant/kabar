// internal/domain/generate/entity.go
package generate

type ArticleGenerationResult struct {
	Article     string `json:"article"`
	Title       string `json:"title"`
	Topic       string `json:"topic"`
	WordCount   int    `json:"wordCount"`
	ImagePrompt string `json:"imagePrompt,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty"`
}

type ImageGenerationResult struct {
	ImageURL string `json:"imageUrl"`
	Prompt   string `json:"prompt"`
}

type ArticleGenerationParams struct {
	Topic             string `json:"topic"`
	ModelID           string `json:"modelId"`
	Tone              string `json:"tone"`
	Length            string `json:"length"`
	Language          string `json:"language"`
	AutoGenerateImage bool   `json:"autoGenerateImage"`
}

type ImageGenerationParams struct {
	Prompt  string `json:"prompt"`
	ModelID string `json:"modelId"`
}
