// internal/interfaces/api/team/handler.go
package team

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/domain/team"
	"seo-backend/internal/helper"
)

type TeamHandler struct {
	teamService team.Service
}

func NewTeamHandler(teamService team.Service) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// ========== Helper Methods ==========

func (h *TeamHandler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (h *TeamHandler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, map[string]string{"error": message}, status)
}

func (h *TeamHandler) handleServiceError(w http.ResponseWriter, err error) {
	errMsg := err.Error()

	switch {
	case errMsg == "access denied":
		h.writeError(w, "Forbidden", http.StatusForbidden)
	case errMsg == "team not found":
		h.writeError(w, "Team not found", http.StatusNotFound)
	case strings.Contains(errMsg, "cannot delete team with active members"):
		h.writeError(w, errMsg, http.StatusBadRequest)
	case strings.Contains(errMsg, "member already in team"):
		h.writeError(w, errMsg, http.StatusConflict)
	case strings.Contains(errMsg, "member not found"):
		h.writeError(w, "Member not found", http.StatusNotFound)
	case strings.Contains(errMsg, "maximum member limit"):
		h.writeError(w, errMsg, http.StatusBadRequest)
	default:
		log.Printf("Unexpected error: %v", err)
		h.writeError(w, "Internal server error", http.StatusInternalServerError)
	}
}

func validateTeamRequest(req team.CreateTeamRequest) error {
	if req.Name == "" {
		return &ValidationError{Message: "team name is required"}
	}
	if len(req.Name) > 100 {
		return &ValidationError{Message: "team name too long (max 100 characters)"}
	}
	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ========== User Context Implementation ==========

type UserContextImpl struct {
	UserID     string
	TeamID     string
	Role       string
	IsAdminVal bool
}

func (u *UserContextImpl) IsAdmin() bool {
	return u.IsAdminVal
}

func (u *UserContextImpl) GetTeamID() string {
	return u.TeamID
}

func (u *UserContextImpl) GetUserID() string {
	return u.UserID
}

func (u *UserContextImpl) GetUserRole() string {
	return u.Role
}

// ========== Team CRUD Handlers ==========

// GetAll returns all teams based on user permissions
func (h *TeamHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := helper.GetUserContext(r)

	filters := team.TeamFilters{
		Status: r.URL.Query().Get("status"),
	}

	teams, err := h.teamService.GetAll(ctx, userCtx, filters)
	if err != nil {
		log.Printf("Failed to fetch teams: %v", err)
		h.writeError(w, "Failed to fetch teams", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, teams, http.StatusOK)
}

// GetByID returns a specific team by ID
func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := helper.GetUserContext(r)
	id := chi.URLParam(r, "id")

	teamData, err := h.teamService.GetByID(ctx, id, userCtx)
	if err != nil {
		log.Printf("Failed to fetch team: %v", err)
		if err.Error() == "access denied" {
			h.writeError(w, "Forbidden", http.StatusForbidden)
		} else {
			h.writeError(w, "Team not found", http.StatusNotFound)
		}
		return
	}

	if teamData == nil {
		h.writeError(w, "Team not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, teamData, http.StatusOK)
}

// Create creates a new team
func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	useCtx := helper.GetUserContext(r)

	var req team.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := validateTeamRequest(req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	teamData, err := h.teamService.Create(ctx, req, useCtx.GetUserID())
	if err != nil {
		log.Printf("Failed to create team: %v", err)
		h.writeError(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, teamData, http.StatusOK)
}

// Update updates an existing team
func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := helper.GetUserContext(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.teamService.Update(ctx, id, updates, userCtx); err != nil {
		log.Printf("Failed to update team: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, map[string]string{"message": "Team updated successfully"}, http.StatusOK)
}

// Delete removes a team
func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := helper.GetUserContext(r)
	id := chi.URLParam(r, "id")

	if err := h.teamService.Delete(ctx, id, userCtx); err != nil {
		log.Printf("Failed to delete team: %v", err)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ========== Member Management Handlers ==========

// GetTeamMembers returns all members of a team
func (h *TeamHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	useCtx := helper.GetUserContext(r)
	teamID := chi.URLParam(r, "id")

	filters := team.MemberFilters{
		Role: r.URL.Query().Get("role"),
	}

	members, err := h.teamService.GetTeamMembers(ctx, teamID, filters, useCtx)
	if err != nil {
		log.Printf("Failed to fetch team members: %v", err)
		h.writeError(w, "Failed to fetch team members", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, members, http.StatusOK)
}

// AddMember adds a user to a team
func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	useCtx := helper.GetUserContext(r)
	teamID := chi.URLParam(r, "id")

	var req team.AddTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		h.writeError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	teamData, err := h.teamService.AddMember(ctx, teamID, req, useCtx)
	if err != nil {
		log.Printf("Failed to add member: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, teamData, http.StatusCreated)
}

// UpdateMemberRole updates a team member's role
func (h *TeamHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := helper.GetUserContext(r)
	teamID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	var role team.TeamMemberRole
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if role == "" {
		h.writeError(w, "role is required", http.StatusBadRequest)
		return
	}

	teamData, err := h.teamService.UpdateMemberRole(ctx, teamID, userID, role, userCtx)
	if err != nil {
		log.Printf("Failed to update member role: %v", err)
		if err.Error() == "member not found" {
			h.writeError(w, "Member not found", http.StatusNotFound)
		} else {
			h.writeError(w, "Failed to update member role", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, teamData, http.StatusOK)
}

// RemoveMember removes a user from a team
func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := helper.GetUserContext(r)
	teamID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	teamData, err := h.teamService.RemoveMember(ctx, teamID, userID, userCtx)
	if err != nil {
		log.Printf("Failed to remove member: %v", err)
		if err != nil {
			h.writeError(w, "Member not found", http.StatusNotFound)
		}
		h.writeError(w, "Failed to remove member", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, teamData, http.StatusOK)
}

// ========== User Teams Handler ==========

// GetUserTeams returns all teams a user belongs to
func (h *TeamHandler) GetUserTeams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userId")

	if userID == "" {
		h.writeError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	teams, err := h.teamService.GetUserTeams(ctx, userID)
	if err != nil {
		log.Printf("Failed to fetch user teams: %v", err)
		h.writeError(w, "Failed to fetch user teams", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, teams, http.StatusOK)
}
