package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/domain/auth"
	authmiddle "seo-backend/internal/middleware"
	"seo-backend/internal/models"
)

type AuthHandler struct {
	service auth.AuthService
}

func NewAuthHandler(service auth.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Helper functions
func (h *AuthHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 403 {object} map[string]string "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: email=%s", req.Email)

	resp, err := h.service.Login(r.Context(), auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Login failed: %v", err)
		// Ubah dari 401 Unauthorized menjadi 403 Forbidden
		h.writeError(w, "Invalid credentials", http.StatusForbidden)
		return
	}

	log.Printf("Login successful: %s", req.Email)

	h.writeJSON(w, models.LoginResponse{
		Token: resp.Token,
		User:  *resp.User,
	}, http.StatusOK)
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} map[string]string "Bad request - missing required fields or invalid password"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), auth.RegisterRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		switch err.Error() {
		case "user already exists":
			h.writeError(w, "User already exists", http.StatusConflict)
		case "email is required", "name is required", "password is required", "password must be at least 6 characters":
			h.writeError(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeError(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("User registered as ADMIN: %s (%s)", user.Email, user.Role)

	h.writeJSON(w, user, http.StatusCreated)
}

// GetMe godoc
// @Summary Get current user
// @Description Get authenticated user information
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} models.User "User details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		} else if err.Error() == "user not found" {
			h.writeError(w, "User not found", http.StatusNotFound)
		} else {
			h.writeError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, user, http.StatusOK)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]string "message: Password changed successfully"
// @Failure 400 {object} map[string]string "Invalid request or new password must be at least 6 characters"
// @Failure 401 {object} map[string]string "Unauthorized or invalid old password"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userID := authmiddle.GetUserID(ctx)

	log.Printf("============USER ID %v", userID)

	err := h.service.ChangePassword(r.Context(), auth.ChangePasswordRequest{
		UserID:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		switch err.Error() {
		case "unauthorized":
			h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		case "user not found":
			h.writeError(w, "User not found", http.StatusNotFound)
		case "invalid old password":
			h.writeError(w, "Invalid old password", http.StatusUnauthorized)
		case "new password must be at least 6 characters":
			h.writeError(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("Failed to change password: %v", err)
			h.writeError(w, "Failed to change password", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, map[string]string{"message": "Password changed successfully"}, http.StatusOK)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset link to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email address"
// @Success 200 {object} map[string]string "message: If your email is registered, you will receive a reset link"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.ForgotPassword(r.Context(), auth.ForgotPasswordRequest{
		Email: req.Email,
	})

	if err != nil {
		log.Printf("Forgot password error: %v", err)
		// Don't reveal error for security
	}

	// Always return same message for security
	h.writeJSON(w, map[string]string{
		"message": "If your email is registered, you will receive a reset link",
	}, http.StatusOK)
}

// Logout godoc
// @Summary Logout user
// @Description Logout current user (client-side token removal)
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "message: Logged out successfully"
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, map[string]string{"message": "Logged out successfully"}, http.StatusOK)
}
