package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
)

type TeamHandler struct{}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{}
}

// ==================== GET ALL TEAMS ====================
func (h *TeamHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)

	status := r.URL.Query().Get("status")

	qb := builder.NewQueryBuilder("teams")

	query := qb.Select(
		"id", "name", "description", "created_by", "created_at", "updated_at",
	).OrderBy("created_at DESC")

	if userRole != "admin" && userRole != "super_admin" {
		if teamID != "" {
			query = query.WhereEq("id", teamID)
		} else {
			query = query.Where("1=0")
		}
	}

	if status != "" {
		query = query.WhereEq("status", status)
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch teams", http.StatusInternalServerError)
		return
	}

	rows, err := database.GetDB().Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch teams", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.Logo, &t.Status,
			&t.MaxMembers, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		members, _ := h.getTeamMembers(t.ID)
		t.Members = members
		teams = append(teams, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// ==================== GET TEAM BY ID ====================
func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)

	id := chi.URLParam(r, "id")

	if userRole != "admin" && userRole != "super_admin" {
		if teamID != "" && id != teamID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	qb := builder.NewQueryBuilder("teams")

	sqlQuery, args, err := qb.Select(
		"id", "name", "description", "logo", "status", "max_members",
		"created_by", "created_at", "updated_at",
	).WhereEq("id", id).Build()

	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch team", http.StatusInternalServerError)
		return
	}

	var team models.Team
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&team.ID, &team.Name, &team.Description, &team.Logo, &team.Status,
		&team.MaxMembers, &team.CreatedBy, &team.CreatedAt, &team.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	members, _ := h.getTeamMembers(team.ID)
	team.Members = members

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

// ==================== CREATE TEAM ====================
func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)

	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if userID == "" {
		userID = "system"
	}

	qb := builder.NewQueryBuilder("teams")

	sqlQuery, args, err := qb.Insert().
		Columns("name", "description", "status", "max_members", "created_by").
		Values(req.Name, req.Description, "active", 10, userID).
		Returning("id").
		Build()

	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	var teamID string
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(&teamID)
	if err != nil {
		log.Printf("Failed to create team: %v", err)
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	// Auto-add creator as team member (tanpa status, invited_by, dll)
	_, err = database.GetDB().Exec(`
		INSERT INTO team_members (team_id, user_id, role, joined_at)
		VALUES ($1, $2, 'manager', $3)
	`, teamID, userID, time.Now())

	if err != nil {
		log.Printf("Failed to add creator as member: %v", err)
	}

	// Get created team
	var team models.Team
	err = database.GetDB().QueryRow(`
		SELECT id, name, description, status, max_members, created_by, created_at, updated_at 
		FROM teams WHERE id = $1
	`, teamID).Scan(&team.ID, &team.Name, &team.Description, &team.Status,
		&team.MaxMembers, &team.CreatedBy, &team.CreatedAt, &team.UpdatedAt)

	if err == nil {
		team.Members = []models.TeamMember{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(team)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      teamID,
		"message": "Team created successfully",
	})
}

// ==================== UPDATE TEAM ====================
func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)

	id := chi.URLParam(r, "id")

	if userRole != "admin" && userRole != "super_admin" {
		if teamID != id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fieldMap := map[string]string{
		"name":        "name",
		"description": "description",
		"logo":        "logo",
		"status":      "status",
		"maxMembers":  "max_members",
	}

	qb := builder.NewQueryBuilder("teams")
	updateBuilder := qb.Update()

	hasUpdates := false
	for key, value := range updates {
		if dbField, ok := fieldMap[key]; ok && value != nil {
			updateBuilder = updateBuilder.Set(dbField, value)
			hasUpdates = true
		}
	}

	if !hasUpdates {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	updateBuilder = updateBuilder.Set("updated_at", time.Now())
	updateBuilder = updateBuilder.WhereEq("id", id)

	sqlQuery, args, err := updateBuilder.Build()
	if err != nil {
		log.Printf("Failed to build update query: %v", err)
		http.Error(w, "Failed to update team", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to update team: %v", err)
		http.Error(w, "Failed to update team", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Team updated successfully"})
}

// ==================== DELETE TEAM ====================
func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)

	id := chi.URLParam(r, "id")

	if userRole != "admin" && userRole != "super_admin" {
		if teamID != id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	// Check member count
	var memberCount int
	err := database.GetDB().QueryRow(`
		SELECT COUNT(*) FROM team_members WHERE team_id = $1
	`, id).Scan(&memberCount)

	if err != nil {
		log.Printf("Failed to check member count: %v", err)
	}

	if memberCount > 0 {
		http.Error(w, "Cannot delete team with active members", http.StatusBadRequest)
		return
	}

	qb := builder.NewQueryBuilder("teams")
	sqlQuery, args, err := qb.Delete().WhereEq("id", id).Build()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		http.Error(w, "Failed to delete team", http.StatusInternalServerError)
		return
	}

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to delete team: %v", err)
		http.Error(w, "Failed to delete team", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ==================== GET TEAM MEMBERS ====================
func (h *TeamHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	members, err := h.getTeamMembers(id)
	if err != nil {
		log.Printf("Failed to fetch team members: %v", err)
		http.Error(w, "Failed to fetch team members", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// ==================== ADD TEAM MEMBER ====================
func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "id")

	var req models.AddTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if member already exists
	var exists bool
	err := database.GetDB().QueryRow(`
		SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)
	`, teamID, req.UserID).Scan(&exists)

	if err != nil {
		log.Printf("Failed to check existing member: %v", err)
	}

	if exists {
		http.Error(w, "Member already in team", http.StatusConflict)
		return
	}

	// Check member limit
	var currentMemberCount int
	var maxMembers int
	err = database.GetDB().QueryRow(`
		SELECT COUNT(*), COALESCE(max_members, 10) 
		FROM team_members tm, teams t
		WHERE tm.team_id = $1 AND t.id = $1
		GROUP BY max_members
	`, teamID).Scan(&currentMemberCount, &maxMembers)

	if err == nil && currentMemberCount >= maxMembers {
		http.Error(w, "Team has reached maximum member limit", http.StatusBadRequest)
		return
	}

	// Insert member (tanpa status, invited_by, updated_at)
	_, err = database.GetDB().Exec(`
		INSERT INTO team_members (team_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)
	`, teamID, req.UserID, req.Role, time.Now())

	if err != nil {
		log.Printf("Failed to add member: %v", err)
		http.Error(w, "Failed to add member", http.StatusInternalServerError)
		return
	}

	team, err := h.GetByIDInternal(teamID)
	if err != nil {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

// ==================== UPDATE TEAM MEMBER ROLE ====================
func (h *TeamHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	var req models.AddTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := database.GetDB().Exec(`
		UPDATE team_members SET role = $1
		WHERE team_id = $2 AND user_id = $3
	`, req.Role, teamID, userID)

	if err != nil {
		log.Printf("Failed to update member role: %v", err)
		http.Error(w, "Failed to update member role", http.StatusInternalServerError)
		return
	}

	team, err := h.GetByIDInternal(teamID)
	if err != nil {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

// ==================== REMOVE TEAM MEMBER ====================
func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	_, err := database.GetDB().Exec(`
		DELETE FROM team_members
		WHERE team_id = $1 AND user_id = $2
	`, teamID, userID)

	if err != nil {
		log.Printf("Failed to remove member: %v", err)
		http.Error(w, "Failed to remove member", http.StatusInternalServerError)
		return
	}

	team, err := h.GetByIDInternal(teamID)
	if err != nil {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

// ==================== GET USER TEAMS ====================
func (h *TeamHandler) GetUserTeams(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	// Raw query dengan JOIN ke team_members (tanpa status)
	query := `
		SELECT t.id, t.name, t.description, t.logo, t.status, t.max_members, 
		       t.created_by, t.created_at, t.updated_at
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY tm.joined_at DESC
	`

	rows, err := database.GetDB().Query(query, userID)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch user teams", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.Logo, &t.Status, &t.MaxMembers,
			&t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		members, _ := h.getTeamMembers(t.ID)
		t.Members = members
		teams = append(teams, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// ==================== INTERNAL HELPER FUNCTIONS ====================

func (h *TeamHandler) getTeamMembers(teamID string) ([]models.TeamMember, error) {
	// Query tanpa status, left_at, invited_by
	query := `
		SELECT id, team_id, user_id, role, joined_at
		FROM team_members
		WHERE team_id = $1
		ORDER BY 
			CASE WHEN role = 'manager' THEN 1 ELSE 2 END,
			joined_at ASC
	`

	rows, err := database.GetDB().Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var m models.TeamMember
		err := rows.Scan(
			&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.JoinedAt,
		)
		if err != nil {
			log.Printf("Scan error in getTeamMembers: %v", err)
			continue
		}
		members = append(members, m)
	}

	return members, nil
}

func (h *TeamHandler) GetByIDInternal(id string) (*models.Team, error) {
	query := `
		SELECT id, name, description, logo, status, max_members, created_by, created_at, updated_at 
		FROM teams WHERE id = $1
	`

	var team models.Team
	err := database.GetDB().QueryRow(query, id).Scan(
		&team.ID, &team.Name, &team.Description, &team.Logo, &team.Status,
		&team.MaxMembers, &team.CreatedBy, &team.CreatedAt, &team.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	members, err := h.getTeamMembers(team.ID)
	if err != nil {
		return nil, err
	}
	team.Members = members

	return &team, nil
}
