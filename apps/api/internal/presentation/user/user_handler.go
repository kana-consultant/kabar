package user

import (
	"encoding/json"
	"log"
	"net/http"
	"seo-backend/common"
	"seo-backend/internal/domain/user"
	"seo-backend/internal/models"

	auth "seo-backend/internal/middleware"

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	switch status {
	case http.StatusBadRequest:
		json.NewEncoder(w).Encode(common.ErrorResponse400{
			Error:  message,
			Status: status,
		})
	case http.StatusUnauthorized:
		json.NewEncoder(w).Encode(common.ErrorResponse401{
			Error:  message,
			Status: status,
		})
	case http.StatusForbidden:
		json.NewEncoder(w).Encode(common.ErrorResponse403{
			Error:  message,
			Status: status,
		})
	case http.StatusNotFound:
		json.NewEncoder(w).Encode(common.ErrorResponse404{
			Error:  message,
			Status: status,
		})
	case http.StatusConflict:
		json.NewEncoder(w).Encode(common.ErrorResponse409{
			Error:  message,
			Status: status,
		})
	default:
		json.NewEncoder(w).Encode(common.ErrorResponse500{
			Error:  message,
			Status: status,
		})
	}
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

// Create godoc
// @Summary Create a new user
// @Description Register a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User registration data"
// @Success 201 {object} models.User "User created"
// @Failure 400 {object} common.ErrorResponse400 "Bad request - invalid email, name, or password"
// @Failure 409 {object} common.ErrorResponse409 "Email already exists"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Router /api/users [post]
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

// GetByID godoc
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User "User details"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "User not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/users/{id} [get]
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

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Description Get the currently authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} models.User "Current user details"
// @Failure 401 {object} common.ErrorResponse401 "Unauthorized"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/users/me [get]
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

// GetAll godoc
// @Summary Get all users
// @Description Get all users with optional filters
// @Tags Users
// @Accept json
// @Produce json
// @Param role query string false "Filter by role" Enums(admin,user,viewer)
// @Param status query string false "Filter by status" Enums(active,inactive,suspended)
// @Param search query string false "Search by name or email"
// @Param limit query int false "Items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {array} models.User "List of users"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/users [get]
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

// Update godoc
// @Summary Update a user
// @Description Update an existing user's information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User update data"
// @Success 200 {object} common.SuccessUpdated "User updated successfully"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request body"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "User not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/users/{id} [put]
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

	h.writeJSON(w, common.SuccessUpdated{
		ID:      id,
		Message: "User updated successfully",
	}, http.StatusOK)
}

// =======================
// UPDATE USER PASSWORD
// =======================

// UpdatePassword godoc
// @Summary Update user password
// @Description Update a user's password (requires old password verification)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdatePasswordRequest true "Password update data"
// @Success 200 {object} common.SuccessUpdated "Password updated successfully"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request or password too short"
// @Failure 403 {object} common.ErrorResponse403 "Access denied or invalid old password"
// @Failure 404 {object} common.ErrorResponse404 "User not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/users/{id}/password [put]

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := h.getUserContext(r)
	id := chi.URLParam(r, "id")

	var req UpdatePasswordRequest

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

	h.writeJSON(w, common.SuccessUpdated{
		ID:      id,
		Message: "Password updated successfully",
	}, http.StatusOK)
}

// =======================
// DELETE USER
// =======================

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by ID (admin only or self-deletion)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} common.ErrorResponse400 "Cannot delete this user"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "User not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/users/{id} [delete]
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
			r.Put("/{id}/password", h.UpdatePassword)
		})
	})
}
