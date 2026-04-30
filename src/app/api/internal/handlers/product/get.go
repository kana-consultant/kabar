package product

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetByID returns a specific product by ID
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	productID := chi.URLParam(r, "id")

	log.Printf("[PRODUCT GET] id=%s user=%s role=%s", productID, userCtx.GetUserID(), userCtx.GetRole())

	product, err := h.productService.GetByID(productID, userCtx)
	if err != nil {
		log.Printf("[PRODUCT GET ERROR] %v", err)
		if err.Error() == "access denied" {
			h.writeError(w, "Forbidden", http.StatusForbidden)
		} else {
			h.writeError(w, "Product not found", http.StatusNotFound)
		}
		return
	}

	h.writeJSON(w, product, http.StatusOK)
}
