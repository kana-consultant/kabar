package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	services "seo-backend/internal/service/generate"
)

type GenerateHandler struct {
	generateService *services.Service
}

func NewGenerateHandler() *GenerateHandler {
	return &GenerateHandler{
		generateService: services.NewService(),
	}
}

func (h *GenerateHandler) GenerateArticle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Topic             string `json:"topic"`
		ModelID           string `json:"modelId"`
		Tone              string `json:"tone"`
		Length            string `json:"length"`
		Language          string `json:"language"`
		AutoGenerateImage bool   `json:"autoGenerateImage"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR decode request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Request: Topic=%s ModelID=%s Tone=%s Length=%s Language=%s",
		req.Topic, req.ModelID, req.Tone, req.Length, req.Language)

	params := services.ArticleGenerationParams{
		Topic:             req.Topic,
		ModelID:           req.ModelID,
		Tone:              req.Tone,
		Length:            req.Length,
		Language:          req.Language,
		AutoGenerateImage: req.AutoGenerateImage,
	}

	result, err := h.generateService.GenerateArticle(params)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *GenerateHandler) GenerateImage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Prompt  string `json:"prompt"`
		ModelID string `json:"modelId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Request: Prompt=%s, ModelID=%s", req.Prompt, req.ModelID)

	params := services.ImageGenerationParams{
		Prompt:  req.Prompt,
		ModelID: req.ModelID,
	}

	result, err := h.generateService.GenerateImage(params)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
