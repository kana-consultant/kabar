// internal/interfaces/api/generate/handler.go
package generate

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/domain/generate"
)

type GenerateHandler struct {
	generateService generate.Service
}

func NewGenerateHandler(generateService generate.Service) *GenerateHandler {
	return &GenerateHandler{
		generateService: generateService,
	}
}

func (h *GenerateHandler) GenerateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Topic             string `json:"topic"`
		ModelID           string `json:"modelId"`
		Tone              string `json:"tone"`
		Length            string `json:"length"`
		Language          string `json:"language"`
		Slug              string `json:"slug,omitempty"`
		ArticleID         string `json:"articleId,omitempty"`
		AutoGenerateImage bool   `json:"autoGenerateImage"`
		ImageModelID      string `json:"imageModelId,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR decode request: %v", err)
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Topic == "" || req.ModelID == "" {
		h.writeError(w, "topic and modelId are required", http.StatusBadRequest)
		return
	}

	// Validate image configuration
	if req.AutoGenerateImage && req.ImageModelID == "" {
		log.Printf("[WARNING] AutoGenerateImage is enabled but ImageModelID is not provided")
		// Allow the request to proceed, but service will handle this case gracefully
	}

	params := generate.ArticleGenerationParams{
		Topic:             req.Topic,
		ModelID:           req.ModelID,
		Tone:              req.Tone,
		Length:            req.Length,
		Language:          req.Language,
		Slug:              req.Slug,
		ArticleID:         req.ArticleID,
		AutoGenerateImage: req.AutoGenerateImage,
		ImageModelID:      req.ImageModelID,
	}

	log.Printf("[INFO] GenerateArticle request: topic=%s, modelId=%s, autoGenerateImage=%v, imageModelId=%s",
		req.Topic, req.ModelID, req.AutoGenerateImage, req.ImageModelID)

	result, err := h.generateService.GenerateArticle(ctx, params)
	if err != nil {
		log.Printf("ERROR: %v", err)
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

func (h *GenerateHandler) GenerateImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Prompt    string `json:"prompt"`
		ModelID   string `json:"modelId"`
		Slug      string `json:"slug,omitempty"`
		ArticleID string `json:"articleId,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request: %v", err)
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" || req.ModelID == "" {
		h.writeError(w, "prompt and modelId are required", http.StatusBadRequest)
		return
	}

	params := generate.ImageGenerationParams{
		Prompt:    req.Prompt,
		ModelID:   req.ModelID,
		Slug:      req.Slug,
		ArticleID: req.ArticleID,
	}

	log.Printf("[INFO] GenerateImage request: prompt=%s, modelId=%s, slug=%s",
		req.Prompt, req.ModelID, req.Slug)

	result, err := h.generateService.GenerateImage(ctx, params)
	if err != nil {
		log.Printf("ERROR: %v", err)
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

func (h *GenerateHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (h *GenerateHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}
