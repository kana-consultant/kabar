package draft

import (
	"strings"

	"golang.org/x/net/html"
)

type SEOScore struct {
	Total       int            `json:"total"`
	Details     map[string]int `json:"details"`
	Suggestions []string       `json:"suggestions"`
}

func CalculateSEOScore(title, content, excerpt, topic string) SEOScore {
	details := map[string]int{}
	suggestions := []string{}
	total := 0

	keyword := strings.ToLower(topic)
	contentLower := strings.ToLower(content)
	titleLower := strings.ToLower(title)

	// 1. Keyword in title (20 pts)
	if strings.Contains(titleLower, keyword) {
		details["keyword_in_title"] = 20
		total += 20
	} else {
		suggestions = append(suggestions, "Add the main keyword to the title")
	}

	// 2. H1 exists (15 pts)
	if strings.Contains(content, "<h1") {
		details["has_h1"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Add an H1 heading to the content")
	}

	// 3. H2 exists (10 pts)
	if strings.Contains(content, "<h2") {
		details["has_h2"] = 10
		total += 10
	} else {
		suggestions = append(suggestions, "Add H2 subheadings to the content")
	}

	// 4. Keyword in first 100 words (15 pts)
	words := strings.Fields(stripHTML(contentLower))
	first100 := strings.Join(words[:min(100, len(words))], " ")
	if strings.Contains(first100, keyword) {
		details["keyword_in_intro"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Use the keyword in the first 100 words")
	}

	// 5. Meta description length 120-160 chars (15 pts)
	excerptLen := len(excerpt)
	if excerptLen >= 120 && excerptLen <= 160 {
		details["meta_description"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Meta description should be 120-160 characters")
	}

	// 6. Content length > 600 words (15 pts)
	wordCount := len(strings.Fields(stripHTML(contentLower)))
	if wordCount >= 600 {
		details["content_length"] = 15
		total += 15
	} else {
		suggestions = append(suggestions, "Content should be at least 600 words")
	}

	// 7. Keyword in content (10 pts)
	if strings.Count(contentLower, keyword) >= 2 {
		details["keyword_density"] = 10
		total += 10
	} else {
		suggestions = append(suggestions, "Use the keyword at least 2 times in the content")
	}

	return SEOScore{
		Total:       total,
		Details:     details,
		Suggestions: suggestions,
	}
}

func stripHTML(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return content
	}
	var buf strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return buf.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}