package builder

const htmlRules = `
==================================================
HTML STRUCTURE & FORMATTING REQUIREMENTS
==================================================

1. ALLOWED TAGS (STRICT):
   <h1> <h2> <h3> <p> <strong> <em>
   <ul> <ol> <li>
   <table> <thead> <tbody> <tr> <th> <td>
   <img> <a>

   FORBIDDEN:
   <div> <span> <section> <br> <style> <script>
   Markdown, backticks, HTML comments, inline styles

2. DOCUMENT STRUCTURE:
   - Exactly ONE <h1> as first element
   - 5-8 <h2> sections
   - <h3> only for subsections (max 3 per H2)
   - Every heading MUST be followed by at least one <p>
   - End with text (<p>), never image or table
   - No empty elements

3. PARAGRAPH RULES:
   - 2-5 sentences per paragraph (50-120 words)
   - Max 3 consecutive paragraphs without break element
   - Mix with lists/tables/images for visual variety

4. HEADING RULES:
   - H1: 5-15 words, include primary keyword naturally
   - H2: 3-10 words, at least 1-2 contain keyword
   - No keyword stuffing, write for humans

5. LIST RULES:
   - Minimum 3 <li> items, maximum 8
   - Introduce lists with a <p>
   - Use <strong> for key terms in <li>

6. TABLE RULES:
   - Only for comparison/structured data
   - Must have <thead> + <th> + <tbody>
   - 2-5 columns, 2-10 rows
   - Introduce tables with a <p>
   - Keep cells concise (1-2 sentences)

7. IMAGE RULES (CRITICAL):
   - EVERY <img> MUST have alt="descriptive text" (5-15 words)
   - Include keyword naturally in alt text
   - NEVER image as first or last element
   - NEVER two consecutive images
   - Max 2-4 images total, evenly distributed
   - Place AFTER referencing paragraph
   - No images in intro (first 150 words) or conclusion

8. LINK RULES:
   - Use <a href="..."> for links
   - Descriptive anchor text (NOT "click here")
   - Include 1-3 links where natural
   - Links should enhance content, not distract

9. TEXT FORMATTING:
   - <strong>: key terms, important points (2-4 per section)
   - <em>: emphasis, technical terms on first use
   - Don't bold/italicize entire sentences
   - Never nest <strong> and <em>

10. STRUCTURE BALANCE:
    - Varied elements: text → list → text → table → text
    - Avoid clustering (no 3+ lists in a row)
    - Section lengths similar (±30%)
    - Smooth transitions between content types

11. VALIDATION:
    - All tags properly closed and nested
    - No duplicate attributes
    - Alt text on all images
    - Table headers on all tables
    - H1 exactly once
    - Start with <h1>, end with <p>
    - No forbidden tags or formatting

12. OUTPUT:
    - Pure HTML only, start directly with <h1>
    - No <html>, <head>, <body> wrappers
    - Valid UTF-8, no hidden characters
    - No markdown, code fences, or XML
`
