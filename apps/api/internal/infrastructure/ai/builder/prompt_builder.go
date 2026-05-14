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
	language := params.Language
	if language == "" {
		language = "English"
	}

	tone := params.Tone
	if tone == "" {
		tone = "professional"
	}

	length := params.Length
	if length == "" {
		length = "medium (800-1200 words)"
	}

	return fmt.Sprintf(`
You are an expert SEO content writer with deep expertise and authority in your field.

Write a high-quality SEO-friendly article in %s about "%s".

Requirements:
- Tone: %s
- Length: %s
- Content must be valid HTML

EEAT Guidelines (strictly follow):
- Experience: Write as someone with firsthand experience on the topic
- Expertise: Use accurate, in-depth information with proper terminology
- Authoritativeness: Reference credible concepts and best practices
- Trustworthiness: Be factual, balanced, and cite reasoning clearly

SEO & Keyword Rules:
- Place the primary keyword in the H1 title
- Use the keyword naturally in the first 100 words
- Include keyword in at least one H2 heading
- Use semantic/related keywords throughout the content
- Add proper meta description in excerpt

Content Structure:
- Use proper HTML heading hierarchy (H1 > H2 > H3)
- Include an introduction, body sections, and conclusion
- Add internal linking suggestions as HTML comments

IMPORTANT: Return ONLY valid JSON, no markdown, no backticks, no explanation.

{
  "title": "string",
  "content": "<h1>...</h1>",
  "excerpt": "string (150-160 chars for meta description)",
  "keywords": ["keyword1", "keyword2"],
  "imagePrompt": "string (prompt for featured image generation)",
  "seoScore": number,
  "readabilityScore": number,
  "wordCount": number
}
`,
		language,
		params.Topic,
		tone,
		length,
	)
}