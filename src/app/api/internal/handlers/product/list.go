package product

import (
	"log"
	"net/http"
	"seo-backend/internal/service/product"
)

// GetAll returns all products based on user permissions
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	filters := product.ProductFilters{
		Platform:   r.URL.Query().Get("platform"),
		Status:     r.URL.Query().Get("status"),
		SyncStatus: r.URL.Query().Get("sync_status"),
	}

	products, err := h.productService.GetAll(userCtx, filters)
	if err != nil {
		log.Printf("Failed to fetch products: %v", err)
		h.writeError(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	log.Printf("Returning %d products", len(products))
	h.writeJSON(w, products, http.StatusOK)
}
