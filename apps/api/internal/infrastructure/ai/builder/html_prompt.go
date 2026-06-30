package builder

const htmlRules = `
==================================================
HTML REQUIREMENTS
==================================================

Generate semantic, clean, and well-structured HTML.

Allowed tags:

<h1>
<h2>
<h3>
<p>
<ul>
<ol>
<li>
<strong>
<em>
<table>
<thead>
<tbody>
<tr>
<th>
<td>
<img>

Document layout:

- Begin with exactly one <h1>.
- Use <h2> for major sections.
- Use <h3> only when a subsection is necessary.
- Every heading should be followed by at least one <p>.
- Keep paragraphs relatively short (2–5 sentences).
- Use <ul> or <ol> whenever listing multiple items.
- Use <table> only when tabular data improves readability.
- Place images immediately after the paragraph that introduces or explains them.
- Never place two consecutive images without meaningful text between them.
- Distribute images naturally throughout the article.
- Avoid placing an image as the first or last element of the article.
- End the article with a text section instead of an image.
- Maintain a balanced flow between headings, paragraphs, lists, tables, and images.
- Do not create empty sections.

Image layout:

- Every <img> must be surrounded by relevant content.
- Images should support the nearby explanation instead of interrupting it.
- Avoid clustering multiple visual elements together.
- Keep image placement consistent with the reading flow.

Do NOT generate:

<style>
<script>
<iframe>
<svg>
<canvas>
markdown
backticks
XML
`
