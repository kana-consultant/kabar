package helper

import "strings"

func GenerateMetaTags(
	title string,
	topic string,
	excerpt string,
	imageURL string,
) string {
	var metaTags strings.Builder

	// Pre-allocate capacity untuk menghindari reallocation
	metaTags.Grow(1024) // Estimasi ukuran

	// Helper untuk menulis string dengan newline
	writeLine := func(s string) {
		metaTags.WriteString(s)
		metaTags.WriteByte('\n')
	}

	// Basic meta tags
	writeLine(`<meta charset="UTF-8">`)
	writeLine(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)

	// Title
	if title != "" {
		writeLine(`<title>` + escapeHTML(title) + `</title>`)
	}

	// Description / Excerpt
	if excerpt != "" {
		writeLine(`<meta name="description" content="` + escapeHTML(excerpt) + `">`)
	} else if topic != "" {
		writeLine(`<meta name="description" content="` + escapeHTML(topic) + `">`)
	}

	// Robots
	writeLine(`<meta name="robots" content="index, follow">`)

	// Open Graph
	if title != "" {
		writeLine(`<meta property="og:title" content="` + escapeHTML(title) + `">`)
	}
	if excerpt != "" {
		writeLine(`<meta property="og:description" content="` + escapeHTML(excerpt) + `">`)
	} else if topic != "" {
		writeLine(`<meta property="og:description" content="` + escapeHTML(topic) + `">`)
	}
	if imageURL != "" {
		writeLine(`<meta property="og:image" content="` + escapeHTML(imageURL) + `">`)
	}
	writeLine(`<meta property="og:type" content="article">`)

	// Twitter Card
	writeLine(`<meta name="twitter:card" content="summary_large_image">`)
	if title != "" {
		writeLine(`<meta name="twitter:title" content="` + escapeHTML(title) + `">`)
	}
	if excerpt != "" {
		writeLine(`<meta name="twitter:description" content="` + escapeHTML(excerpt) + `">`)
	}
	if imageURL != "" {
		writeLine(`<meta name="twitter:image" content="` + escapeHTML(imageURL) + `">`)
	}

	// Keywords
	if topic != "" {
		writeLine(`<meta name="keywords" content="` + escapeHTML(topic) + `">`)
	}

	return metaTags.String()
}

// escapeHTML mengamankan string dari XSS
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return strings.ReplaceAll(s, "'", "&#39;")
}

// Helper function to determine publish status
func DeterminePublishStatus(someFailed bool, allFailed bool, err error) string {
	if err != nil || allFailed {
		return "failed"
	}
	if someFailed {
		return "partial"
	}
	return "published"
}

// Helper function to get status message
func GetStatusMessage(someFailed bool, allFailed bool) string {
	if allFailed {
		return "All products failed to publish"
	}
	if someFailed {
		return "Some products failed to publish"
	}
	return "All products published successfully"
}
