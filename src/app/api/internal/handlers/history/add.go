package history

import (
	"encoding/json"
	"log"
	"net/http"
	"seo-backend/internal/service/history"
)

// AddToHistory adds a new history record
func (h *HistoryHandler) AddToHistory(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	var req struct {
		Title          string   `json:"title"`
		Topic          string   `json:"topic"`
		Content        string   `json:"content"`
		ImageURL       *string  `json:"imageUrl,omitempty"`
		TargetProducts []string `json:"targetProducts"`
		Status         string   `json:"status"`
		Action         string   `json:"action"`
		ErrorMessage   *string  `json:"errorMessage,omitempty"`
		ScheduledFor   *string  `json:"scheduledFor,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	createReq := history.CreateHistoryRequest{
		Title:          req.Title,
		Topic:          req.Topic,
		Content:        req.Content,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
		Status:         req.Status,
		Action:         req.Action,
		ErrorMessage:   req.ErrorMessage,
		ScheduledFor:   req.ScheduledFor,
		CreatedBy:      userCtx.GetUserID(),
		TeamID:         userCtx.GetTeamID(),
	}

	historyID, err := h.historyService.Create(createReq)
	if err != nil {
		log.Printf("Failed to add to history: %v", err)
		h.writeError(w, "Failed to add to history", http.StatusInternalServerError)
		return
	}

	log.Printf("History added with ID: %s", historyID)
	h.writeJSON(w, map[string]string{
		"id":      historyID,
		"message": "Added to history successfully",
	}, http.StatusCreated)
}
