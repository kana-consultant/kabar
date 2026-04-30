package product

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Delete removes a product
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userCtx := h.getUserContext(r)

	if err := h.productService.Delete(id, userCtx); err != nil {
		log.Printf("Failed to delete product: %v", err)
		if err.Error() == "access denied" {
			h.writeError(w, "Forbidden", http.StatusForbidden)
		} else if err.Error() == "product not found" {
			h.writeError(w, "Product not found", http.StatusNotFound)
		} else {
			h.writeError(w, "Failed to delete product", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
