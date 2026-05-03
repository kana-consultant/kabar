package apikey

import (
	"encoding/json"
	"net/http"

	"seo-backend/internal/database"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := database.GetDB().Exec(`
		DELETE FROM api_keys WHERE id = $1
	`, id)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to check result", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "deleted successfully",
		"id":      id,
	})
}
