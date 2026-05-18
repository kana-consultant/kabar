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
	return fmt.Sprintf(`
TOPIC:
"%s"

CORE RULES:
- Write naturally like a human writer
- Ensure the article is informative, structured, and readable
- Avoid repetitive sentences and generic AI phrases
- Do not use placeholders or dummy text
- Keep explanations relevant to the topic
- Use a clear logical flow between sections
- Maintain factual consistency

CONTENT REQUIREMENTS:
- Use valid HTML only
- Use proper heading hierarchy:
  <h1>, <h2>, <h3>
- Include:
  - Introduction
  - Main discussion sections
  - Conclusion
- Use paragraphs, lists, and tables when relevant
- Do not include markdown
- Do not wrap output with backticks

SEO REQUIREMENTS:
- Include the main topic naturally in:
  - Title
  - Introduction
  - At least one subheading
- Use semantic and related keywords naturally
- Optimize readability
- Avoid keyword stuffing

OUTPUT RULES:
- Return ONLY valid JSON
- No explanation
- No additional text outside JSON
- Ensure JSON is parseable

RESPONSE FORMAT:
{
  "title": "string",
  "slug": "string",
  "excerpt": "string",
  "content": "<h1>...</h1>",
  "keywords": ["string"],
  "imagePrompt": "string",
  "wordCount": number
}
`,
		params.Topic,
	)
}
