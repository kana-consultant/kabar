package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
	services "seo-backend/internal/service/product"
)

type ProductHandler struct {
	productService *services.Service
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: services.NewService(),
	}
}

// Helper: Get user context from request
func (h *ProductHandler) getUserContext(r *http.Request) models.UserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}

// Helper: Write JSON response
func (h *ProductHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper: Write error response
func (h *ProductHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// Helper: Write error response with success flag
func (h *ProductHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}

// Helper: Handle service errors
func (h *ProductHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "access denied":
		h.writeError(w, "Forbidden", http.StatusForbidden)
	case "product not found":
		h.writeError(w, "Product not found", http.StatusNotFound)
	default:
		log.Printf("Unexpected error: %v", err)
		h.writeError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Helper: Validate product request
func validateProductRequest(req models.CreateProductRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Platform == "" {
		return fmt.Errorf("platform is required")
	}
	if req.APIEndpoint == "" {
		return fmt.Errorf("api endpoint is required")
	}
	return nil
}
