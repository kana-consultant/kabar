// internal/infrastructure/ai/builder/prompt_builder.go
package builder

import (
	"fmt"

	"seo-backend/internal/domain/generate"
)

type PromptBuilder struct{}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

func (b *PromptBuilder) BuildArticlePrompt(params generate.ArticleGenerationParams) string {
	imageSection := noImageRules

	if params.AutoGenerateImage {
		imageSection = imageRules
	}

	return fmt.Sprintf(`
      ==================================================
      PRIMARY OBJECTIVE
      ==================================================

      Generate one complete SEO article that fully answers the user's search intent.

      Return ONLY valid JSON.

      ==================================================
      TOPIC
      ==================================================

      %s

      %s

      %s

      %s

      %s

      %s
`,
		params.Topic,
		articleRules,
		htmlRules,
		seoRules,
		imageSection,
		outputRules,
	)
}
