package generate

// Generation parameters
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

// Results
type ArticleResult struct {
	Title            string   `json:"title"`
	Content          string   `json:"content"`
	Excerpt          string   `json:"excerpt"`
	Keywords         []string `json:"keywords"`
	ImagePrompt      string   `json:"imagePrompt"`
	ImageURL         string   `json:"imageUrl"`
	WordCount        int      `json:"wordCount"`
	ReadabilityScore int      `json:"readabilityScore"`
	SeoScore         int      `json:"seoScore"`
}

type ImageResult struct {
	ImageURL    string `json:"imageUrl"`
	Prompt      string `json:"prompt"`
	GeneratedAt string `json:"generatedAt"`
	Model       string `json:"model"`
}

// Internal types
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
}
