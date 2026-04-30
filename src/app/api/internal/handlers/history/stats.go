package history

import (
	"log"
	"net/http"
)

// GetStats returns statistics about history records
func (h *HistoryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	stats, err := h.historyService.GetStats(userCtx)
	if err != nil {
		log.Printf("Failed to fetch stats: %v", err)
		h.writeError(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, stats, http.StatusOK)
}
