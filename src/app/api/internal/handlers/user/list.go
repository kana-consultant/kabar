package user

import (
	"log"
	"net/http"
	"seo-backend/internal/service/user"
)

// GetAll returns all users based on user permissions
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userCtx := h.getUserContext(r)

	filters := user.UserFilters{
		Role:   r.URL.Query().Get("role"),
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
	}

	users, err := h.userService.GetAll(userCtx, filters)
	if err != nil {
		log.Printf("Failed to fetch users: %v", err)
		h.writeError(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully fetched %d users", len(users))
	h.writeJSON(w, users, http.StatusOK)
}
