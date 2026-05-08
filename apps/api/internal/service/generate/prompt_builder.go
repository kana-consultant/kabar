package generate

import "fmt"

func buildArticlePrompt(params ArticleGenerationParams) string {
	// Set defaults if empty
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
Write a high-quality SEO-friendly article in %s about "%s".

Requirements:
- Tone: %s
- Length: %s
- Content must be valid HTML
- SEO optimized with proper heading structure
- Include internal linking suggestions

Return ONLY valid JSON.

{
"title":"string",
"content":"<h1>...</h1>",
"excerpt":"string",
"seoScore":number,
"readabilityScore":number,
"wordCount":number
}
`,
		language,
		params.Topic,
		tone,
		length,
	)
}
