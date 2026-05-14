// internal/interfaces/api/aimodel/handler.go
package aimodel

import (
	"encoding/json"
	"log"
	"net/http"

	aimodel "seo-backend/internal/domain/model"
	auth "seo-backend/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

type AIModelHandler struct {
	aimodelService aimodel.Service
}

func NewAIModelHandler(aimodelService aimodel.Service) *AIModelHandler {
	return &AIModelHandler{
		aimodelService: aimodelService,
	}
}

// GetAll returns all active AI models
func (h *AIModelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	models, err := h.aimodelService.GetAllModels(ctx)
	if err != nil {
		log.Printf("Failed to fetch models: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, models, http.StatusOK)
}

// GetAllWithStatus returns all models with API key status
func (h *AIModelHandler) GetAllWithStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)

	models, err := h.aimodelService.GetAllModelsWithStatus(ctx, userRole)
	if err != nil {
		log.Printf("Failed to fetch models with status: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, models, http.StatusOK)
}

// GetByID returns a specific model by ID
func (h *AIModelHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	model, err := h.aimodelService.GetModelByID(ctx, id)
	if err != nil {
		log.Printf("Failed to fetch model: %v", err)
		h.writeError(w, "Failed to fetch model", http.StatusInternalServerError)
		return
	}

	if model == nil {
		h.writeError(w, "Model not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, model, http.StatusOK)
}

// GetByProvider returns models by provider ID
func (h *AIModelHandler) GetByProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	providerID := chi.URLParam(r, "providerId")

	models, err := h.aimodelService.GetModelsByProvider(ctx, providerID)
	if err != nil {
		log.Printf("Failed to fetch models by provider: %v", err)
		h.writeError(w, "Failed to fetch models", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, models, http.StatusOK)
}

// Helper methods
func (h *AIModelHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (h *AIModelHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}
