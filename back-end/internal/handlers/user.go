package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"seo-backend/internal/builder"
	"seo-backend/internal/database"
	"seo-backend/internal/middleware/auth"
	"seo-backend/internal/models"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// ==================== GET ALL USERS ====================
// ==================== GET ALL USERS ====================
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	// Build query dengan builder
	qb := builder.NewQueryBuilder("users")

	query := qb.Select(
		"id", "email", "name", "role", "avatar", "status", "last_active", "created_at", "updated_at",
	).OrderBy("created_at DESC")

	// Logic filter berdasarkan role
	switch userRole {
	case "super_admin":
		// Super admin: melihat SEMUA user
		// Tidak ada filter tambahan
		log.Printf("Super admin - melihat semua user")

	case "admin":
		// Admin: hanya melihat user dalam team yang sama
		if teamID != "" {
			query = query.Where(`
				id IN (
					SELECT user_id FROM team_members WHERE team_id = $1
				)
			`, teamID)
			log.Printf("Admin - melihat user dalam team: %s", teamID)
		} else {
			// Admin tanpa team hanya melihat dirinya sendiri
			query = query.WhereEq("id", userID)
			log.Printf("Admin tanpa team - hanya melihat diri sendiri")
		}

	default:
		// Role lain (manager, editor, viewer): hanya melihat diri sendiri
		query = query.WhereEq("id", userID)
		log.Printf("User role %s - hanya melihat diri sendiri", userRole)
	}

	sqlQuery, args, err := query.Build()
	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	log.Printf("Query: %s | Args: %v", sqlQuery, args)

	rows, err := database.GetDB().Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Email, &u.Name, &u.Role, &u.Avatar, &u.Status, &u.LastActive, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		users = append(users, u)
	}

	log.Printf("Successfully fetched %d users", len(users))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// ==================== GET USER BY ID ====================
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole := auth.GetUserRole(ctx)
	teamID := auth.GetTeamID(ctx)
	userID := auth.GetUserID(ctx)

	id := chi.URLParam(r, "id")

	// Cek authorization
	if userRole != "admin" && userRole != "super_admin" && id != userID {
		var isSameTeam bool
		err := database.GetDB().QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM team_members 
				WHERE team_id = $1 AND user_id = $2
			)
		`, teamID, id).Scan(&isSameTeam)

		if err != nil || !isSameTeam {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	qb := builder.NewQueryBuilder("users")

	sqlQuery, args, err := qb.Select(
		"id", "email", "name", "role", "avatar", "status", "last_active", "created_at", "updated_at",
	).WhereEq("id", id).Build()

	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar, &user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ==================== GET CURRENT USER ====================
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(ctx)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	qb := builder.NewQueryBuilder("users")

	sqlQuery, args, err := qb.Select(
		"id", "email", "name", "role", "avatar", "status", "last_active", "created_at", "updated_at",
	).WhereEq("id", userID).Build()

	if err != nil {
		log.Printf("Failed to build query: %v", err)
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = database.GetDB().QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar, &user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ==================== CREATE USER ====================
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	role := models.RoleViewer
	if req.Role != "" {
		role = models.UserRole(req.Role)
	}

	var user models.User
	query := `INSERT INTO users (email, name, password_hash, role, status) 
	          VALUES ($1, $2, $3, $4, 'active') 
	          RETURNING id, email, name, role, avatar, status, last_active, created_at, updated_at`
	err = database.GetDB().QueryRow(query, req.Email, req.Name, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar, &user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// ==================== UPDATE USER ====================
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Map field names dari JSON ke database
	fieldMap := map[string]string{
		"name":   "name",
		"email":  "email",
		"role":   "role",
		"status": "status",
	}

	qb := builder.NewQueryBuilder("users")
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

	updateBuilder = updateBuilder.Set("updated_at", "NOW()")
	updateBuilder = updateBuilder.WhereEq("id", id)

	sqlQuery, args, err := updateBuilder.Build()
	if err != nil {
		log.Printf("Failed to build update query: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	log.Printf("Update query: %s | Args: %v", sqlQuery, args)

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

// ==================== DELETE USER ====================
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	qb := builder.NewQueryBuilder("users")
	sqlQuery, args, err := qb.Delete().WhereEq("id", id).Build()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	log.Printf("Delete query: %s | Args: %v", sqlQuery, args)

	_, err = database.GetDB().Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
