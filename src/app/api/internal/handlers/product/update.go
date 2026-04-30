package product

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Update updates an existing product
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userCtx := h.getUserContext(r)

	if err := h.productService.Update(id, updates, userCtx); err != nil {
		log.Printf("Failed to update product: %v", err)
		if err.Error() == "access denied" {
			h.writeError(w, "Forbidden", http.StatusForbidden)
		} else if err.Error() == "product not found" {
			h.writeError(w, "Product not found", http.StatusNotFound)
		} else {
			h.writeError(w, "Failed to update product", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, map[string]string{"message": "Product updated successfully"}, http.StatusOK)
}
