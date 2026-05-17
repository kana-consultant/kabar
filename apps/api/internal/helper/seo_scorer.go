package helper

import (
	"math"
	"regexp"
	"strings"
)

func CalculateSEOScoreSimple(title, content, excerpt, topic string) int {
	total := 0
	keyword := strings.ToLower(topic)
	contentLower := strings.ToLower(content)
	titleLower := strings.ToLower(title)

	if strings.Contains(titleLower, keyword) {
		total += 20
	}
	if strings.Contains(content, "<h1") {
		total += 15
	}
	if strings.Contains(content, "<h2") {
		total += 10
	}

	words := strings.Fields(stripHTMLSimple(contentLower))
	first100 := strings.Join(words[:minInt(100, len(words))], " ")
	if strings.Contains(first100, keyword) {
		total += 15
	}

	excerptLen := len(excerpt)
	if excerptLen >= 120 && excerptLen <= 160 {
		total += 15
	}

	wordCount := len(strings.Fields(stripHTMLSimple(contentLower)))
	if wordCount >= 600 {
		total += 15
	}

	if strings.Count(contentLower, keyword) >= 2 {
		total += 10
	}

	_ = math.Round(0)
	return total
}

func stripHTMLSimple(content string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(content, " ")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}