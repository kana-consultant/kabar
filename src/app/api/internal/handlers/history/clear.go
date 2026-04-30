package history

import (
	"log"
	"net/http"
)

// ClearAll clears all history records (admin only)
func (h *HistoryHandler) ClearAll(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	if err := h.historyService.ClearAll(userCtx); err != nil {
		log.Printf("Failed to clear history: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{"message": "All history cleared successfully"}, http.StatusOK)
}
