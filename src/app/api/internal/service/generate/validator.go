package generate

import "fmt"

func validateArticleParams(params ArticleGenerationParams) error {
	if params.Topic == "" || params.ModelID == "" {
		return fmt.Errorf("topic and modelId are required")
	}
	return nil
}

func validateImageParams(params ImageGenerationParams) error {
	if params.Prompt == "" || params.ModelID == "" {
		return fmt.Errorf("prompt and modelId are required")
	}
	return nil
}
