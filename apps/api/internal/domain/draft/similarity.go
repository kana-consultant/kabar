package draft

import (
	"math"
	"regexp"
	"strings"
)

func ComputeTFIDF(docs []string) []map[string]float64 {
	// Tokenize semua dokumen
	tokenized := make([][]string, len(docs))
	for i, doc := range docs {
		tokenized[i] = tokenize(doc)
	}

	// Hitung IDF
	idf := computeIDF(tokenized)

	// Hitung TF-IDF vector per dokumen
	vectors := make([]map[string]float64, len(docs))
	for i, tokens := range tokenized {
		vectors[i] = computeTF(tokens, idf)
	}

	return vectors
}

func CosineSimilarity(a, b map[string]float64) float64 {
	var dot, normA, normB float64

	for k, va := range a {
		if vb, ok := b[k]; ok {
			dot += va * vb
		}
		normA += va * va
	}

	for _, vb := range b {
		normB += vb * vb
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func tokenize(text string) []string {
	// Strip HTML
	re := regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, " ")

	// Lowercase dan split
	text = strings.ToLower(text)
	words := strings.Fields(text)

	// Filter kata pendek
	var tokens []string
	for _, w := range words {
		w = strings.Trim(w, ".,!?;:\"'()[]{}") 
		if len(w) > 2 {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

func computeIDF(docs [][]string) map[string]float64 {
	df := make(map[string]int)
	for _, tokens := range docs {
		seen := make(map[string]bool)
		for _, t := range tokens {
			if !seen[t] {
				df[t]++
				seen[t] = true
			}
		}
	}

	idf := make(map[string]float64)
	n := float64(len(docs))
	for term, count := range df {
		idf[term] = math.Log(n / float64(count))
	}
	return idf
}

func computeTF(tokens []string, idf map[string]float64) map[string]float64 {
	tf := make(map[string]int)
	for _, t := range tokens {
		tf[t]++
	}

	tfidf := make(map[string]float64)
	total := float64(len(tokens))
	for term, count := range tf {
		tfidf[term] = (float64(count) / total) * idf[term]
	}
	return tfidf
}