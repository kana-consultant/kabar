package builder

const htmlRules = `
==================================================
HTML REQUIREMENTS
==================================================

Use valid HTML only.

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

Do NOT generate:

<style>
<script>
iframe
svg
canvas
markdown
backticks
XML
`
