package history

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetByID returns a specific history record by ID
func (h *HistoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	history, err := h.historyService.GetByID(id, userCtx)
	if err != nil {
		log.Printf("Failed to fetch history: %v", err)
		h.handleServiceError(w, err)
		return
	}

	if history == nil {
		h.writeError(w, "History not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, history, http.StatusOK)
}
