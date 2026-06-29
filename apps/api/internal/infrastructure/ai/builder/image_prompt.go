package builder

const imageRules = `
==================================================
IMAGE PLACEHOLDERS
==================================================

Insert one image placeholder every 500–700 words.

Maximum four image placeholders.

Format:

<img prompt="..." style="max-width:100%;height:auto;border-radius:8px;" />

Prompt rules:

- 20–40 words
- Describe only factual visualizations
- Use charts, graphs, timelines, flowcharts, comparison tables, or infographics
- Images must support the surrounding content

Do NOT request:

- Stock photos
- Portraits
- Logos
- Decorative illustrations
- Header banners
`

const noImageRules = `
==================================================
IMAGE REQUIREMENTS
==================================================

Do NOT generate any <img> tags.

Present information using HTML elements only.

Use:

- <table>
- <ul>
- <ol>

when appropriate.
`
