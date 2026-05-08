package generate

type Client struct {
	service *Service
}

func NewClient() *Client {
	return &Client{
		service: NewService(),
	}
}

// GenerateArticleWithDefaults generates an article with default values
func (c *Client) GenerateArticleWithDefaults(topic, modelID string) (*ArticleResult, error) {
	return c.service.GenerateArticle(ArticleGenerationParams{
		Topic:    topic,
		ModelID:  modelID,
		Tone:     "professional",
		Length:   "medium",
		Language: "English",
	})
}

// GenerateImageWithPrompt generates an image with a simple prompt
func (c *Client) GenerateImageWithPrompt(prompt, modelID string) (*ImageResult, error) {
	return c.service.GenerateImage(ImageGenerationParams{
		Prompt:  prompt,
		ModelID: modelID,
	})
}
