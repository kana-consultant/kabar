package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"seo-backend/internal/domain/auth"
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

// Login handles user login
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
		h.writeError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("Login successful: %s", req.Email)

	h.writeJSON(w, models.LoginResponse{
		Token: resp.Token,
		User:  *resp.User,
	}, http.StatusOK)
}

// Register handles user registration
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

// GetMe returns current authenticated user
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

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

// ForgotPassword handles forgot password request
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

// Logout handles logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, map[string]string{"message": "Logged out successfully"}, http.StatusOK)
}
