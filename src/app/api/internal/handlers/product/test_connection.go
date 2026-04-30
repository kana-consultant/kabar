package product

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TestConnection tests the connection to a product's API endpoint
func (h *ProductHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.productService.TestConnection(id)
	if err != nil {
		log.Printf("Connection test failed: %v", err)
		h.writeErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}
