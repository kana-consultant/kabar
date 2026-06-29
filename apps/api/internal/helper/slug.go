package helper

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateSlug membuat slug dari string dengan berbagai opsi
func GenerateSlug(text string) string {
	return GenerateSlugWithOptions(text, SlugOptions{
		Lowercase:       true,
		RemoveStopwords: false,
		MaxLength:       100,
		Separator:       "-",
		KeepNumbers:     true,
	})
}

// SlugOptions adalah konfigurasi untuk generate slug
type SlugOptions struct {
	Lowercase       bool   // Ubah ke lowercase (default: true)
	RemoveStopwords bool   // Hapus stopwords (default: false)
	MaxLength       int    // Maksimal panjang slug (default: 100)
	Separator       string // Separator karakter (default: "-")
	KeepNumbers     bool   // Pertahankan angka (default: true)
	KeepDots        bool   // Pertahankan titik (default: false)
	KeepUnderscores bool   // Pertahankan underscore (default: false)
}

// GenerateSlugWithOptions membuat slug dengan opsi kustom
func GenerateSlugWithOptions(text string, opts SlugOptions) string {
	if text == "" {
		return ""
	}

	// Set default options
	if opts.Separator == "" {
		opts.Separator = "-"
	}
	if opts.MaxLength == 0 {
		opts.MaxLength = 100
	}

	// Trim spasi
	text = strings.TrimSpace(text)

	// Convert ke lowercase jika diperlukan
	if opts.Lowercase {
		text = strings.ToLower(text)
	}

	// Hapus stopwords jika diperlukan
	if opts.RemoveStopwords {
		text = removeStopwords(text)
	}

	// Replace special characters
	text = replaceSpecialChars(text)

	// Hanya pertahankan karakter yang diizinkan
	allowedChars := "a-z0-9"
	if opts.KeepDots {
		allowedChars += "."
	}
	if opts.KeepUnderscores {
		allowedChars += "_"
	}

	reg := regexp.MustCompile(fmt.Sprintf(`[^%s]`, allowedChars))
	text = reg.ReplaceAllString(text, " ")

	// Hapus multiple spaces
	spaceReg := regexp.MustCompile(`\s+`)
	text = spaceReg.ReplaceAllString(text, " ")

	// Trim spasi
	text = strings.TrimSpace(text)

	// Replace spasi dengan separator
	text = strings.ReplaceAll(text, " ", opts.Separator)

	// Hapus multiple separators
	sepReg := regexp.MustCompile(fmt.Sprintf(`%s+`, regexp.QuoteMeta(opts.Separator)))
	text = sepReg.ReplaceAllString(text, opts.Separator)

	// Trim separator di awal dan akhir
	text = strings.Trim(text, opts.Separator)

	// Batasi panjang
	if len(text) > opts.MaxLength {
		text = text[:opts.MaxLength]
		// Hapus separator di akhir jika ada
		text = strings.TrimSuffix(text, opts.Separator)
	}

	return text
}

// GenerateSlugFromTitle membuat slug dari title dengan timestamp (untuk uniqueness)
func GenerateSlugFromTitle(title string) string {
	slug := GenerateSlug(title)

	// Tambahkan timestamp pendek untuk uniqueness
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", slug, timestamp)
}

// GenerateSlugWithID membuat slug dengan ID (untuk uniqueness)
func GenerateSlugWithID(title string, id string) string {
	slug := GenerateSlug(title)

	// Ambil beberapa karakter terakhir dari ID
	suffix := id
	if len(id) > 8 {
		suffix = id[len(id)-8:]
	}

	return fmt.Sprintf("%s-%s", slug, suffix)
}

// GenerateSlugUnique membuat slug unique dengan counter
func GenerateSlugUnique(title string, existingSlugs []string) string {
	baseSlug := GenerateSlug(title)

	// Cek apakah slug sudah ada
	exists := false
	for _, s := range existingSlugs {
		if s == baseSlug {
			exists = true
			break
		}
	}

	if !exists {
		return baseSlug
	}

	// Tambahkan counter
	counter := 1
	for {
		newSlug := fmt.Sprintf("%s-%d", baseSlug, counter)
		exists := false
		for _, s := range existingSlugs {
			if s == newSlug {
				exists = true
				break
			}
		}
		if !exists {
			return newSlug
		}
		counter++
	}
}

// Helper: replace special characters
func replaceSpecialChars(text string) string {
	// Map karakter khusus ke ASCII
	replacements := map[string]string{
		"à": "a", "á": "a", "â": "a", "ã": "a", "ä": "a", "å": "a",
		"è": "e", "é": "e", "ê": "e", "ë": "e",
		"ì": "i", "í": "i", "î": "i", "ï": "i",
		"ò": "o", "ó": "o", "ô": "o", "õ": "o", "ö": "o", "ø": "o",
		"ù": "u", "ú": "u", "û": "u", "ü": "u",
		"ý": "y", "ÿ": "y",
		"ñ": "n", "ç": "c",
		"ß": "ss",
		"æ": "ae", "œ": "oe",
		"&": "and",
		"@": "at",
		"+": "plus",
		"%": "percent",
		"#": "number",
		"¢": "cent",
		"€": "euro",
		"£": "pound",
		"¥": "yen",
		"©": "copyright",
		"®": "registered",
		"™": "tm",
	}

	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	return text
}

// Helper: remove stopwords
func removeStopwords(text string) string {
	// Stopwords bahasa Indonesia dan Inggris
	stopwords := map[string]bool{
		// Indonesia
		"yang": true, "dan": true, "di": true, "ke": true, "dari": true,
		"ini": true, "itu": true, "adalah": true, "untuk": true, "dengan": true,
		"pada": true, "dalam": true, "oleh": true, "sebagai": true, "atau": true,
		"saya": true, "kamu": true, "kami": true, "mereka": true, "dia": true,
		"akan": true, "telah": true, "bisa": true, "dapat": true, "juga": true,
		"sangat": true, "cukup": true, "sekali": true, "saja": true, "pun": true,
		"lagi": true, "masih": true, "sudah": true, "belum": true, "pernah": true,

		// Inggris
		"the": true, "be": true, "to": true, "of": true, "and": true,
		"a": true, "in": true, "that": true, "have": true, "i": true,
		"it": true, "for": true, "not": true, "on": true, "with": true,
		"he": true, "as": true, "you": true, "do": true, "at": true,
		"this": true, "but": true, "his": true, "by": true, "from": true,
		"they": true, "we": true, "say": true, "her": true, "she": true,
		"or": true, "an": true, "will": true, "my": true, "one": true,
		"all": true, "would": true, "there": true, "their": true, "what": true,
		"so": true, "up": true, "out": true, "if": true, "about": true,
		"who": true, "get": true, "which": true, "go": true, "me": true,
		"when": true, "make": true, "can": true, "like": true, "time": true,
		"no": true, "just": true, "him": true, "know": true, "take": true,
		"people": true, "into": true, "year": true, "your": true, "good": true,
		"some": true, "could": true, "them": true, "see": true, "other": true,
		"than": true, "then": true, "now": true, "look": true, "only": true,
		"come": true, "its": true, "over": true, "think": true, "also": true,
		"back": true, "after": true, "use": true, "two": true, "how": true,
		"our": true, "work": true, "first": true, "well": true, "way": true,
		"even": true, "new": true, "want": true, "because": true, "any": true,
		"these": true, "give": true, "day": true, "most": true, "us": true,
	}

	words := strings.Fields(text)
	var result []string

	for _, word := range words {
		// Hapus punctuation
		word = strings.Trim(word, ".,!?;:()[]{}'\"")
		if !stopwords[strings.ToLower(word)] {
			result = append(result, word)
		}
	}

	return strings.Join(result, " ")
}

// GenerateSlugFromTopic membuat slug dari topic dengan format SEO-friendly
func GenerateSlugFromTopic(topic string) string {
	// Generate base slug
	slug := GenerateSlug(topic)

	// Jika terlalu pendek, tambahkan kata "article"
	if len(slug) < 10 {
		return fmt.Sprintf("artikel-%s", slug)
	}

	return slug
}

// GenerateSlugFromTitleWithDate membuat slug dengan tanggal
func GenerateSlugFromTitleWithDate(title string) string {
	slug := GenerateSlug(title)
	date := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s-%s", slug, date)
}

// GenerateSlugFromTitleWithCategory membuat slug dengan kategori
func GenerateSlugFromTitleWithCategory(title string, category string) string {
	slugTitle := GenerateSlug(title)
	slugCategory := GenerateSlug(category)
	return fmt.Sprintf("%s/%s", slugCategory, slugTitle)
}

// IsValidSlug memvalidasi apakah slug valid
func IsValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	// Cek karakter yang diizinkan
	valid := regexp.MustCompile(`^[a-z0-9\-_\.]+$`)
	return valid.MatchString(slug)
}

// NormalizeSlug membersihkan slug yang tidak valid
func NormalizeSlug(slug string) string {
	// Hapus karakter yang tidak valid
	reg := regexp.MustCompile(`[^a-z0-9\-_\.]`)
	normalized := reg.ReplaceAllString(strings.ToLower(slug), "")

	// Hapus multiple separators
	sepReg := regexp.MustCompile(`[-_\.]{2,}`)
	normalized = sepReg.ReplaceAllString(normalized, "-")

	// Trim separators
	normalized = strings.Trim(normalized, "-_.")

	return normalized
}
