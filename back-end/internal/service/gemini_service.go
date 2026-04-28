package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

type GeminiService struct {
	client *genai.Client
	apiKey string // <-- tambahkan ini
	ctx    context.Context
}

type ArticleRequest struct {
	Topic    string   `json:"topic"`
	Tone     string   `json:"tone"`     // professional, casual, friendly, formal
	Language string   `json:"language"` // id, en
	Keywords []string `json:"keywords"`
	Length   string   `json:"length"` // short, medium, long
}

type ArticleResponse struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Excerpt     string   `json:"excerpt"`
	Keywords    []string `json:"keywords"`
	SeoScore    int      `json:"seoScore"`
	WordCount   int      `json:"wordCount"`
	ImagePrompt string   `json:"imagePrompt"`
}

func NewGeminiService(apiKey string) (*GeminiService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiService{
		client: client,
		apiKey: apiKey, // <-- simpan apiKey
		ctx:    ctx,
	}, nil
}

func (s *GeminiService) GenerateArticle(req ArticleRequest) (*ArticleResponse, error) {
	// Build prompt berdasarkan parameter
	prompt := s.buildArticlePrompt(req)

	// Panggil Gemini API
	result, err := s.client.Models.GenerateContent(
		s.ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Printf("Gemini API error: %v", err)
		return nil, fmt.Errorf("failed to generate article: %w", err)
	}

	// Parse response
	content := result.Text()
	if content == "" {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	// Extract JSON dari response
	jsonStr := extractJSON(content)

	var response ArticleResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return nil, fmt.Errorf("failed to parse article response")
	}

	// Hitung word count
	response.WordCount = len(strings.Fields(response.Content))

	// Generate image prompt dari title
	response.ImagePrompt = fmt.Sprintf("Professional illustration for article: %s, modern, clean design, high quality", response.Title)

	return &response, nil
}

func (s *GeminiService) buildArticlePrompt(req ArticleRequest) string {
	toneDesc := map[string]string{
		"professional": "professional and formal",
		"casual":       "casual and easy to understand",
		"friendly":     "friendly and engaging",
		"formal":       "formal and structured",
	}

	tone := toneDesc[req.Tone]
	if tone == "" {
		tone = toneDesc["professional"]
	}

	lengthDesc := map[string]string{
		"short":  "300-400 words",
		"medium": "600-800 words",
		"long":   "1000-1500 words",
	}
	length := lengthDesc[req.Length]
	if length == "" {
		length = lengthDesc["medium"]
	}

	language := req.Language
	if language == "" {
		language = "id"
	}

	languageName := "Indonesian"
	if language == "en" {
		languageName = "English"
	}

	keywords := strings.Join(req.Keywords, ", ")
	if keywords == "" {
		keywords = req.Topic
	}

	return fmt.Sprintf(`You are a professional SEO content writer. Write a high-quality article in %s language.

TOPIC: %s
TONE: %s
LENGTH: %s
KEYWORDS to include: %s

Return ONLY valid JSON format with these fields:
{
  "title": "catchy article title",
  "content": "full article content with HTML formatting (use <h2>, <p>, <ul>/<li> for structure)",
  "excerpt": "short summary under 160 characters",
  "keywords": ["array", "of", "keywords"],
  "seoScore": (integer score 0-100 based on SEO best practices)
}

The article should be well-structured, informative, and engaging. Use proper HTML tags for formatting.
DO NOT include any text outside the JSON.`,
		languageName, req.Topic, tone, length, keywords)
}

func (s *GeminiService) GenerateImage(prompt string, outputPath string) (string, error) {
	enhancedPrompt := fmt.Sprintf(`%s

Style: professional, high quality, modern, clean design
Lighting: soft natural lighting
Composition: centered, balanced
Quality: 4K, photorealistic

NO text, NO watermark, NO logo, NO people faces if not requested`, prompt)

	resp, err := s.client.Models.GenerateContent(
		s.ctx,
		"gemini-3-pro-image-preview",
		genai.Text(enhancedPrompt),
		nil,
	)
	if err != nil {
		log.Printf("Image generation error: %v", err)
		return "", fmt.Errorf("failed to generate image: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	var imageData []byte
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			imageData = part.InlineData.Data
			break
		}
	}

	if len(imageData) == 0 {
		return "", fmt.Errorf("no image data in response")
	}

	// Check if data is base64 encoded
	imageDataStr := string(imageData)
	if !strings.HasPrefix(imageDataStr, "data:image") {
		decoded, err := base64.StdEncoding.DecodeString(imageDataStr)
		if err == nil {
			imageData = decoded
		}
	}

	if outputPath == "" {
		outputPath = fmt.Sprintf("output/generated-%d.png", time.Now().Unix())
	}

	if err := os.MkdirAll("output", 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(outputPath, imageData, 0644); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (s *GeminiService) EditImage(imagePath string, prompt string, outputPath string) (string, error) {
	// Baca image file
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// Tentukan MIME type
	mimeType := "image/png"
	if strings.HasSuffix(imagePath, ".jpg") || strings.HasSuffix(imagePath, ".jpeg") {
		mimeType = "image/jpeg"
	}

	// Encode ke base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Prompt editing
	editPrompt := fmt.Sprintf(`Edit this image: %s
Keep the original quality and style.
Maintain the same dimensions and resolution.`, prompt)

	// Gunakan API key dari service (perlu ditambahkan field apiKey ke struct)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-3-pro-image-preview:generateContent?key=%s", s.apiKey)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"inlineData": map[string]interface{}{
							"mimeType": mimeType,
							"data":     imageBase64,
						},
					},
					{"text": editPrompt},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(requestBody)

	httpClient := &http.Client{Timeout: 60 * time.Second}
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse response
	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					InlineData struct {
						Data string `json:"data"`
					} `json:"inlineData"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	editedBase64 := geminiResp.Candidates[0].Content.Parts[0].InlineData.Data
	editedData, err := base64.StdEncoding.DecodeString(editedBase64)
	if err != nil {
		return "", err
	}

	if outputPath == "" {
		outputPath = fmt.Sprintf("output/edited-%d.png", time.Now().Unix())
	}

	if err := os.MkdirAll("output", 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(outputPath, editedData, 0644); err != nil {
		return "", err
	}

	return outputPath, nil
}

// Helper function to extract JSON from text
func extractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}
	return "{}"
}
