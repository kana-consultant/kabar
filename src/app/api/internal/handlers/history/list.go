package history

import (
	"log"
	"net/http"
	"seo-backend/internal/service/history"
)

// GetAll returns all history records based on user permissions
func (h *HistoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	filter := history.HistoryFilter{
		Status: r.URL.Query().Get("status"),
		Action: r.URL.Query().Get("action"),
		Search: r.URL.Query().Get("search"),
	}

	// Parse pagination
	filter.Limit, filter.Offset = parsePagination(r)

	histories, err := h.historyService.GetAll(userCtx, filter)
	if err != nil {
		log.Printf("Failed to fetch history: %v", err)
		h.writeError(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully fetched %d history records", len(histories))
	h.writeJSON(w, histories, http.StatusOK)
}
