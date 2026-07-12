// internal/interfaces/api/team/handler.go
package team

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"seo-backend/common"
	"seo-backend/internal/domain/team"

	auth "seo-backend/internal/middleware"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	switch status {
	case http.StatusBadRequest:
		json.NewEncoder(w).Encode(common.ErrorResponse400{
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

// GetAll godoc
// @Summary Get all teams
// @Description Get all teams based on user permissions with optional status filter
// @Tags Teams
// @Accept json
// @Produce json
// @Param status query string false "Filter by team status" Enums(active, inactive, archived)
// @Success 200 {array} team.Team "List of teams"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams [get]
func (h *TeamHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)

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

// GetByID godoc
// @Summary Get team by ID
// @Description Get a specific team by its ID
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} team.Team "Team details"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "Team not found"
// @Security BearerAuth
// @Router /api/teams/{id} [get]
func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
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

// Create godoc
// @Summary Create a new team
// @Description Create a new team with the current user as admin
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body team.CreateTeamRequest true "Team creation request"
// @Success 200 {object} team.Team "Team created"
// @Failure 400 {object} common.ErrorResponse400 "Bad request - invalid team name"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams [post]
func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	useCtx := auth.GetUserContext(r)

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

// Update godoc
// @Summary Update a team
// @Description Update an existing team by ID
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param request body map[string]interface{} true "Team update fields"
// @Success 200 {object} common.SuccessMessage "Team updated successfully"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "Team not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/{id} [put]
func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
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

	h.writeJSON(w, common.SuccessMessage{Message: "Team updated successfully"}, http.StatusOK)
}

// Delete godoc
// @Summary Delete a team
// @Description Delete a team by ID (must be empty)
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 204 "No Content"
// @Failure 400 {object} common.ErrorResponse400 "Cannot delete team with active members"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "Team not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/{id} [delete]
func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
	id := chi.URLParam(r, "id")

	if err := h.teamService.Delete(ctx, id, userCtx); err != nil {
		log.Printf("Failed to delete team: %v", err)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ========== Member Management Handlers ==========

// GetTeamMembers godoc
// @Summary Get team members
// @Description Get all members of a team with optional role filter
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param role query string false "Filter by role" Enums(admin, member, viewer)
// @Success 200 {array} team.TeamMember "List of team members"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/{id}/members [get]
func (h *TeamHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	useCtx := auth.GetUserContext(r)
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

// AddMember godoc
// @Summary Add member to team
// @Description Add a user to a team
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param request body team.AddTeamMemberRequest true "Add member request"
// @Success 201 {object} team.Team "Team with updated members"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request or user_id required"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 409 {object} common.ErrorResponse409 "Member already in team"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/{id}/members [post]
func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	useCtx := auth.GetUserContext(r)
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

// UpdateMemberRole godoc
// @Summary Update member role
// @Description Update a team member's role
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param userId path string true "User ID"
// @Param request body team.TeamMemberRole true "New role"
// @Success 200 {object} team.Team "Team with updated members"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request or role required"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "Member not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/{teamId}/members/{userId}/role [put]
func (h *TeamHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
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

// RemoveMember godoc
// @Summary Remove member from team
// @Description Remove a user from a team
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param userId path string true "User ID"
// @Success 200 {object} team.Team "Team with updated members"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 404 {object} common.ErrorResponse404 "Member not found"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/{id}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := auth.GetUserContext(r)
	teamID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	teamData, err := h.teamService.RemoveMember(ctx, teamID, userID, userCtx)
	if err != nil {
		log.Printf("Failed to remove member: %v", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeJSON(w, teamData, http.StatusOK)
}

// ========== User Teams Handler ==========

// GetUserTeams godoc
// @Summary Get user's teams
// @Description Get all teams a user belongs to
// @Tags Teams
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} team.Team "List of teams"
// @Failure 400 {object} common.ErrorResponse400 "user_id is required"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/user/{userId} [get]
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

// InviteMember godoc
// @Summary Invite member to team
// @Description Send invitation to a user to join a team
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Param request body team.InviteTeamMemberRequest true "Invite request"
// @Success 201 {object} team.TeamInvite "Invitation created"
// @Failure 400 {object} common.ErrorResponse400 "Bad request"
// @Failure 403 {object} common.ErrorResponse403 "Access denied"
// @Failure 500 {object} common.ErrorResponse500 "Internal server error"
// @Security BearerAuth
// @Router /api/teams/{teamId}/invites [post]
func (h *TeamHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	var req team.InviteTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	teamID := chi.URLParam(r, "teamId")
	userCtx := auth.GetUserContext(r)

	invite, err := h.teamService.InviteMember(r.Context(), teamID, req, userCtx)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, invite, http.StatusCreated)
}

// VerificationInvite godoc
// @Summary Verify invitation token
// @Description Verify a team invitation token
// @Tags Teams
// @Accept json
// @Produce json
// @Param token path string true "Invitation token"
// @Success 200 {object} map[string]bool "isTrue: true/false"
// @Failure 400 {object} common.ErrorResponse400 "Invalid token"
// @Router /api/verify-invite/{token} [get]
func (h *TeamHandler) VerificationInvite(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	log.Printf("TOKEN ====== %v", token)

	isTrue, err := h.teamService.VerificationInvite(r.Context(), token)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, isTrue, http.StatusOK)
}

// AcceptInvite godoc
// @Summary Accept team invitation
// @Description Accept a team invitation and join the team
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body team.UserInvitedCreate true "Accept invite request with token"
// @Success 200 {object} map[string]interface{} "success: true, message, teamId, team"
// @Failure 400 {object} common.ErrorResponse400 "Invalid request or token"
// @Security BearerAuth
// @Router /api/accept-invite [post]
func (h *TeamHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	var req team.UserInvitedCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userCtx := auth.GetUserContext(r)

	teamData, err := h.teamService.AcceptInvite(r.Context(), userCtx, req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "Successfully joined team",
		"teamId":  teamData.ID,
		"team":    teamData,
	}, http.StatusOK)
}
