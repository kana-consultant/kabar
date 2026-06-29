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

Generate a complete, accurate, and SEO-optimized article that satisfies the user's search intent.

Return ONLY valid JSON.

==================================================
USER TOPIC
==================================================

%s

==================================================
ARTICLE REQUIREMENTS
==================================================

%s

==================================================
HTML REQUIREMENTS
==================================================

%s

==================================================
SEO REQUIREMENTS
==================================================

%s

%s

==================================================
OUTPUT FORMAT
==================================================

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
