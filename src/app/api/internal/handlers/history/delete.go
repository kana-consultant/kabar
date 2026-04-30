package history

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Delete removes a history record
func (h *HistoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	if err := h.historyService.Delete(id, userCtx); err != nil {
		log.Printf("Failed to delete history: %v", err)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
