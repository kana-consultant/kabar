package product

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/models"
)

// Create creates a new product
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Create product called")

	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateProductRequest(req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCtx := h.getUserContext(r)

	product, err := h.productService.Create(req, userCtx)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		h.writeError(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	log.Printf("Product created: %s", product.Name)
	h.writeJSON(w, product, http.StatusCreated)
}
