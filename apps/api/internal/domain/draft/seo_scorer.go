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

	titleLower := strings.ToLower(title)
	plainContent := stripHTML(content)
	plainLower := strings.ToLower(plainContent)

	// Normalisasi semua keyword ke lowercase
	normalizedKeywords := make([]string, len(keywords))
	for i, k := range keywords {
		normalizedKeywords[i] = strings.ToLower(k)
	}

	log.Printf("[SEO] keywords=%v | title=%q", normalizedKeywords, title)
	log.Printf("[SEO] plain_content_preview=%q", plainLower[:min(200, len(plainLower))])

	// Helper: cek apakah salah satu keyword ada di teks
	containsAnyKeyword := func(text string) bool {
		for _, kw := range normalizedKeywords {
			if strings.Contains(text, kw) {
				return true
			}
		}
		return false
	}

	// Helper: hitung total kemunculan semua keyword
	countAllKeywords := func(text string) int {
		count := 0
		for _, kw := range normalizedKeywords {
			count += strings.Count(text, kw)
		}
		return count
	}

	// 1. Keyword in title (20)
	if containsAnyKeyword(titleLower) {
		details["keyword_in_title"] = 20
		total += 20
		log.Printf("[SEO] ✅ keyword_in_title +20")
	} else {
		suggestions = append(suggestions, "Add the main keyword to the title")
		log.Printf("[SEO] ❌ keyword_in_title | keywords=%v not found in title=%q", normalizedKeywords, titleLower)
	}

	// 2. H1 tag (15)
	if strings.Contains(content, "<h1") {
		details["has_h1"] = 15
		total += 15
		log.Printf("[SEO] ✅ has_h1 +15")
	} else {
		suggestions = append(suggestions, "Add an H1 heading to the content")
		log.Printf("[SEO] ❌ has_h1 | no <h1> found in content")
	}

	// 3. H2 tag (10)
	if strings.Contains(content, "<h2") {
		details["has_h2"] = 10
		total += 10
		log.Printf("[SEO] ✅ has_h2 +10")
	} else {
		suggestions = append(suggestions, "Add H2 subheadings to the content")
		log.Printf("[SEO] ❌ has_h2 | no <h2> found in content")
	}

	// 4. Keyword in intro (15)
	words := strings.Fields(plainLower)
	first100 := strings.Join(words[:min(100, len(words))], " ")
	log.Printf("[SEO] total_words=%d | first_100_words=%q", len(words), first100[:min(100, len(first100))])
	if containsAnyKeyword(first100) {
		details["keyword_in_intro"] = 15
		total += 15
		log.Printf("[SEO] ✅ keyword_in_intro +15")
	} else {
		suggestions = append(suggestions, "Use the keyword in the first 100 words")
		log.Printf("[SEO] ❌ keyword_in_intro | keywords=%v not found in first 100 words", normalizedKeywords)
	}

	// 5. Meta description length (15)
	excerptLen := len(strings.TrimSpace(excerpt))
	log.Printf("[SEO] excerpt_length=%d | excerpt=%q", excerptLen, excerpt)
	if excerptLen >= 120 && excerptLen <= 160 {
		details["meta_description"] = 15
		total += 15
		log.Printf("[SEO] ✅ meta_description +15 | length=%d", excerptLen)
	} else {
		suggestions = append(suggestions, "Meta description should be 120-160 characters")
		log.Printf("[SEO] ❌ meta_description | length=%d (expected 120-160)", excerptLen)
	}

	// 6. Content length (15)
	wordCount := len(strings.Fields(plainLower))
	log.Printf("[SEO] word_count=%d", wordCount)
	if wordCount >= 600 {
		details["content_length"] = 15
		total += 15
		log.Printf("[SEO] ✅ content_length +15 | words=%d", wordCount)
	} else {
		suggestions = append(suggestions, "Content should be at least 600 words")
		log.Printf("[SEO] ❌ content_length | words=%d (expected >=600)", wordCount)
	}

	// 7. Keyword density (10)
	keywordCount := countAllKeywords(plainLower)
	log.Printf("[SEO] keyword_count=%d | keywords=%v", keywordCount, normalizedKeywords)
	if keywordCount >= 2 {
		details["keyword_density"] = 10
		total += 10
		log.Printf("[SEO] ✅ keyword_density +10 | count=%d", keywordCount)
	} else {
		suggestions = append(suggestions, "Use the keyword at least 2 times in the content")
		log.Printf("[SEO] ❌ keyword_density | count=%d (expected >=2)", keywordCount)
	}

	log.Printf("[SEO] total_score=%d | details=%v | suggestions=%v", total, details, suggestions)

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
