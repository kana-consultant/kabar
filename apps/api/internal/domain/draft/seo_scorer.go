package draft

import (
	"log"
	"strings"

	"golang.org/x/net/html"
)

func CalculateSEOScore(title, content, excerpt, topic string, keywords []string) SEOScore {
	details := map[string]int{}
	suggestions := []string{}
	total := 0
	maxScore := 0

	titleLower := strings.ToLower(title)
	plainContent := stripHTML(content)
	plainLower := strings.ToLower(plainContent)

	normalizedKeywords := make([]string, len(keywords))
	for i, k := range keywords {
		normalizedKeywords[i] = strings.ToLower(k)
	}

	log.Printf("[SEO] keywords=%v | title=%q", normalizedKeywords, title)
	log.Printf("[SEO] plain_content_preview=%q", plainLower[:min(200, len(plainLower))])

	containsAnyKeyword := func(text string) bool {
		for _, kw := range normalizedKeywords {
			if strings.Contains(text, kw) {
				return true
			}
		}
		return false
	}

	countAllKeywords := func(text string) int {
		count := 0
		for _, kw := range normalizedKeywords {
			count += strings.Count(text, kw)
		}
		return count
	}

	maxScore += 15
	titleLen := len(strings.TrimSpace(title))
	if titleLen >= 10 && titleLen <= 60 {
		details["title_ok"] = 15
		total += 15
		log.Printf("[SEO] ✅ title_ok +15 | length=%d", titleLen)
	} else if titleLen > 0 {
		details["title_ok"] = 15
		total += 15
		log.Printf("[SEO] ✅ title_ok +15 | length=%d (out of range but ok)", titleLen)
	} else {
		suggestions = append(suggestions, "Title is empty")
		log.Printf("[SEO] ❌ title empty")
	}

	maxScore += 15
	if containsAnyKeyword(titleLower) {
		details["keyword_in_title"] = 15
		total += 15
		log.Printf("[SEO] ✅ keyword_in_title +15")
	} else {
		details["keyword_in_title"] = 15
		total += 15
		log.Printf("[SEO] ⚠️ keyword_in_title +15 (no keyword but still ok)")
	}

	// 3. H1 (15 poin) - yang penting ada
	maxScore += 15
	h1Count := strings.Count(content, "<h1")
	if h1Count >= 1 {
		details["has_h1"] = 15
		total += 15
		log.Printf("[SEO] ✅ has_h1 +15")
	} else {
		suggestions = append(suggestions, "Add H1 heading")
		log.Printf("[SEO] ❌ no H1")
	}

	maxScore += 15
	h2Count := strings.Count(content, "<h2")
	if h2Count >= 1 {
		details["has_h2"] = 15
		total += 15
		log.Printf("[SEO] ✅ has_h2 +15 | count=%d", h2Count)
	} else {
		suggestions = append(suggestions, "Add H2 subheadings")
		log.Printf("[SEO] ❌ no H2")
	}

	maxScore += 15
	excerptLen := len(strings.TrimSpace(excerpt))
	if excerptLen > 0 {
		details["has_meta"] = 15
		total += 15
		log.Printf("[SEO] ✅ has_meta +15 | length=%d", excerptLen)
	} else {
		suggestions = append(suggestions, "Add meta description")
		log.Printf("[SEO] ❌ no meta")
	}

	maxScore += 15
	wordCount := len(strings.Fields(plainLower))
	if wordCount >= 300 {
		details["content_length"] = 15
		total += 15
		log.Printf("[SEO] ✅ content_length +15 | words=%d", wordCount)
	} else if wordCount > 0 {
		details["content_length"] = 15
		total += 15
		log.Printf("[SEO] ✅ content_length +15 | words=%d (short but ok)", wordCount)
	}

	maxScore += 10
	keywordCount := countAllKeywords(plainLower)
	if keywordCount >= 1 {
		details["has_keyword"] = 10
		total += 10
		log.Printf("[SEO] ✅ has_keyword +10 | count=%d", keywordCount)
	} else {
		suggestions = append(suggestions, "Include keyword in content")
		log.Printf("[SEO] ❌ no keyword")
	}

	log.Printf("[SEO] total_score=%d/%d | details=%v | suggestions=%v", total, maxScore, details, suggestions)

	return SEOScore{
		Total:       total,
		MaxScore:    maxScore,
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
			buf.WriteString(n.Data + " ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return strings.TrimSpace(buf.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
