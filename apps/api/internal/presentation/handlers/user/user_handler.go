package user

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/domain/user"
	"seo-backend/internal/models"
	auth "seo-backend/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	service user.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// =======================
// HELPERS
// =======================

// getUserContext extracts user context from request
func (h *UserHandler) getUserContext(r *http.Request) models.UserContext {
	ctx := r.Context()
	return &models.SimpleUserContext{
		UserID: auth.GetUserID(ctx),
		TeamID: auth.GetTeamID(ctx),
		Role:   auth.GetUserRole(ctx),
	}
}

// writeJSON writes JSON response
func (h *UserHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// writeError writes error response
func (h *UserHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// handleServiceError handles service layer errors
func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "access denied":
		h.writeError(w, "Forbidden", http.StatusForbidden)
	case "user not found":
		h.writeError(w, "User not found", http.StatusNotFound)
	case "unauthorized", "unauthorized: no user context":
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
	case "email already exists":
		h.writeError(w, "Email already exists", http.StatusConflict)
	case "cannot delete this user":
		h.writeError(w, "Cannot delete this user", http.StatusBadRequest)
	case "invalid email format", "email is required":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "name is required", "name must be at least 2 characters":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "password is required", "password must be at least 8 characters":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	case "invalid role", "invalid status":
		h.writeError(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("Unexpected error: %v", err)
		h.writeError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// =======================
// CREATE USER
// =======================
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body object true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Prepare request
	createReq := user.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Role:     req.Role,
	}

	// Call service
	result, err := h.service.Create(ctx, createReq)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, result, http.StatusCreated)
}

// =======================
// GET USER BY ID
// =======================
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	result, err := h.service.GetByID(ctx, id, userCtx)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	if result == nil {
		h.writeError(w, "User not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// =======================
// GET CURRENT USER
// =======================
// @Summary Get current authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)

	result, err := h.service.GetCurrentUser(ctx, userCtx)
	if err != nil {
		log.Printf("Failed to fetch current user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, result, http.StatusOK)
}

// =======================
// GET ALL USERS
// =======================
// @Summary Get all users
// @Tags users
// @Produce json
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Param search query string false "Search by name or email"
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)

	users, err := h.service.GetAll(ctx, userCtx)
	if err != nil {
		log.Printf("Failed to fetch users: %v", err)
		h.handleServiceError(w, err)
		return
	}

	log.Printf("Successfully fetched %d users", len(users))
	h.writeJSON(w, users, http.StatusOK)
}

// =======================
// UPDATE USER
// =======================
// @Summary Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body object true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Build update request
	updateReq := user.UpdateUserRequest{}

	if name, ok := updates["name"].(string); ok {
		updateReq.Name = &name
	}
	if email, ok := updates["email"].(string); ok {
		updateReq.Email = &email
	}
	if role, ok := updates["role"].(string); ok {
		updateReq.Role = &role
	}
	if status, ok := updates["status"].(string); ok {
		updateReq.Status = &status
	}
	if avatar, ok := updates["avatar"].(string); ok {
		updateReq.Avatar = &avatar
	}

	if err := h.service.Update(ctx, id, updateReq, userCtx); err != nil {
		log.Printf("Failed to update user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "User updated successfully",
	}, http.StatusOK)
}

// =======================
// UPDATE USER PASSWORD
// =======================
// @Summary Update user password
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body object true "Password data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/password [post]
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OldPassword == "" {
		h.writeError(w, "Old password is required", http.StatusBadRequest)
		return
	}
	if req.NewPassword == "" {
		h.writeError(w, "New password is required", http.StatusBadRequest)
		return
	}
	if len(req.NewPassword) < 8 {
		h.writeError(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePassword(ctx, id, req.OldPassword, req.NewPassword, userCtx); err != nil {
		log.Printf("Failed to update password: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{
		"id":      id,
		"message": "Password updated successfully",
	}, http.StatusOK)
}

// =======================
// DELETE USER
// =======================
// @Summary Delete a user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(ctx, id, userCtx); err != nil {
		log.Printf("Failed to delete user: %v", err)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =======================
// REGISTER ROUTES
// =======================
// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		// Public routes (create user - registration)
		r.Post("/", h.Create)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			// r.Use(middleware.AuthRequired) // if not using global auth

			r.Get("/me", h.GetCurrentUser)
			r.Get("/", h.GetAll)
			r.Get("/{id}", h.GetByID)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Post("/{id}/password", h.UpdatePassword)
		})
	})
}
