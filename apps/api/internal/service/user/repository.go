package user

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"seo-backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, name, role, avatar, status, last_active, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetAll(query string, args []interface{}) ([]models.User, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

func (r *Repository) Create(req CreateUserRequest, passwordHash []byte) (*models.User, error) {
	role := req.Role
	if role == "" {
		role = models.RoleViewer
	}

	query := `
		INSERT INTO users (email, name, password_hash, role, status) 
		VALUES ($1, $2, $3, $4, 'active') 
		RETURNING id, email, name, role, avatar, status, last_active, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(query, req.Email, req.Name, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Avatar,
		&user.Status, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *Repository) Update(id string, updates map[string]interface{}) error {
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIndex)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) Delete(id string) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (r *Repository) scanUsers(rows *sql.Rows) ([]models.User, error) {
	var users []models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Email, &u.Name, &u.Role, &u.Avatar,
			&u.Status, &u.LastActive, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}
