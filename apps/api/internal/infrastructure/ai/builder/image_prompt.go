package builder

const imageRules = `
==================================================
IMAGE PLACEHOLDERS
==================================================

Insert one image placeholder every 500–700 words.

Maximum four image placeholders.

Format:
<img prompt="..." style="max-width:100%;height:auto;border-radius:8px;" />

Prompt requirements:
- 20–40 words.
- Describe only the image to generate.
- The prompt must match the surrounding content.
- Include important details from the content when relevant.
- Be clear, concise, and specific.
- Do not explain the article.
- Do not repeat prompts.

Do NOT request:
- Logos
- Watermarks
- Text overlays
- Unrelated decorative elements
`

const noImageRules = `
==================================================
IMAGE REQUIREMENTS
==================================================

Do NOT generate any <img> tags.

Present information using HTML elements only.

Use:
- <table> for numeric data
- <ul> / <ol> for lists
`
