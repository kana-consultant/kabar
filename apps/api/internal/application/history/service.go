package history

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"seo-backend/internal/domain/history"
	"seo-backend/internal/domain/paginate"
	"seo-backend/internal/helper"
	historyBuilder "seo-backend/internal/infrastructure/db/query_builder"
	"seo-backend/internal/models"
)

// Service handles all history business logic
type Service struct {
	repo         history.HistoryRepository
	queryBuilder historyBuilder.QueryBuilder
}

// NewService creates a new history service
func NewService(repo history.HistoryRepository, queryBuilder historyBuilder.QueryBuilder) *Service {
	return &Service{
		repo:         repo,
		queryBuilder: queryBuilder,
	}
}

// =======================
// CREATE
// =======================

// Create creates a new history record
func (s *Service) Create(ctx context.Context, req history.CreateHistoryRequest) (string, error) {
	// Validation
	if err := s.validateCreateRequest(req); err != nil {
		return "", err
	}

	id := uuid.New().String()
	now := helper.ParseWIBTime(time.Now().Format(time.RFC3339))

	data := &history.History{
		ID:             id,
		Title:          req.Title,
		Topic:          req.Topic,
		Content:        req.Content,
		ImageURL:       req.ImageURL,
		TargetProducts: req.TargetProducts,
		Status:         "pending",
		Action:         "created",
		ErrorMessage:   nil,
		PublishedAt:    now,
		ScheduledFor:   req.ScheduledFor,
		CreatedBy:      &req.CreatedBy,
		TeamID:         &req.TeamID,
		CreatedAt:      now,
	}

	return s.repo.Create(ctx, data)
}

// validateCreateRequest validates create request
func (s *Service) validateCreateRequest(req history.CreateHistoryRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	if req.CreatedBy == "" {
		return fmt.Errorf("createdBy is required")
	}
	if req.TeamID == "" {
		return fmt.Errorf("teamId is required")
	}
	return nil
}

// =======================
// READ
// =======================

// GetAll retrieves all history records based on user context
func (s *Service) GetAll(ctx context.Context, userCtx *models.UserContext, filters history.HistoryFilter) (*paginate.PaginatedResult[history.History], error) {

	return s.repo.GetAll(ctx, *userCtx, filters)
}

// GetByID retrieves a history record by ID
func (s *Service) GetByID(ctx context.Context, id string) (*history.History, error) {
	if id == "" {
		return nil, fmt.Errorf("history id is required")
	}
	return s.repo.GetByID(ctx, id)
}

// GetWithFilters retrieves history with filters
func (s *Service) GetWithFilters(ctx context.Context, filters history.HistoryFilter, user models.UserContext) ([]history.History, int, error) {

	// Build query from filters
	query, args := s.queryBuilder.BuildListQuery(user, filters)

	// Get total count
	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count history: %w", err)
	}

	// Get records
	records, err := s.repo.GetAllWithQuery(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get history: %w", err)
	}

	return records, total, nil
}

// GetByTeamID retrieves history by team ID
func (s *Service) GetByTeamID(ctx context.Context, teamID string) ([]history.History, error) {
	if teamID == "" {
		return nil, fmt.Errorf("team id is required")
	}
	return s.repo.GetByTeamID(ctx, teamID)
}

// GetByCreatedBy retrieves history by creator
func (s *Service) GetByCreatedBy(ctx context.Context, createdBy string) ([]history.History, error) {
	if createdBy == "" {
		return nil, fmt.Errorf("createdBy is required")
	}
	return s.repo.GetByCreatedBy(ctx, createdBy)
}

// GetByStatus retrieves history by status
func (s *Service) GetByStatus(ctx context.Context, status string) ([]history.History, error) {
	if status == "" {
		return nil, fmt.Errorf("status is required")
	}
	return s.repo.GetByStatus(ctx, status)
}

// =======================
// UPDATE
// =======================

// UpdateRequest represents request to update history
type UpdateHistoryRequest struct {
	Title          *string    `json:"title,omitempty"`
	Topic          *string    `json:"topic,omitempty"`
	Content        *string    `json:"content,omitempty"`
	ImageURL       *string    `json:"imageUrl,omitempty"`
	TargetProducts []string   `json:"targetProducts,omitempty"`
	Status         *string    `json:"status,omitempty"`
	Action         *string    `json:"action,omitempty"`
	ErrorMessage   *string    `json:"errorMessage,omitempty"`
	ScheduledFor   *time.Time `json:"scheduledFor,omitempty"`
}

// Update updates an existing history record
func (s *Service) Update(ctx context.Context, id string, req UpdateHistoryRequest) error {
	if id == "" {
		return fmt.Errorf("history id is required")
	}

	// Check if exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("history not found")
	}

	// Build updates
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Topic != nil {
		updates["topic"] = *req.Topic
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.TargetProducts != nil {
		updates["target_products"] = req.TargetProducts
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Action != nil {
		updates["action"] = *req.Action
	}
	if req.ErrorMessage != nil {
		updates["error_message"] = *req.ErrorMessage
	}
	if req.ScheduledFor != nil {
		updates["scheduled_for"] = *req.ScheduledFor
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	return s.repo.Update(ctx, id, updates)
}

// UpdateStatus updates the status of a history record
func (s *Service) UpdateStatus(ctx context.Context, id string, status string, errorMessage *string) error {
	if id == "" {
		return fmt.Errorf("history id is required")
	}
	if status == "" {
		return fmt.Errorf("status is required")
	}

	updates := map[string]interface{}{
		"status": status,
	}
	if errorMessage != nil {
		updates["error_message"] = *errorMessage
	}

	return s.repo.Update(ctx, id, updates)
}

// =======================
// DELETE
// =======================

// Delete removes a history record
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("history id is required")
	}

	// Check if exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("history not found")
	}

	return s.repo.Delete(ctx, id)
}

// DeleteByTeamID removes all history records for a team
func (s *Service) DeleteByTeamID(ctx context.Context, teamID string) error {
	if teamID == "" {
		return fmt.Errorf("team id is required")
	}
	return s.repo.DeleteByTeamID(ctx, teamID)
}

// DeleteByStatus removes history records by status
func (s *Service) DeleteByStatus(ctx context.Context, status string) error {
	if status == "" {
		return fmt.Errorf("status is required")
	}
	return s.repo.DeleteByStatus(ctx, status)
}

// =======================
// STATISTICS
// =======================

// GetStats retrieves statistics for history records
func (s *Service) GetStats(ctx context.Context, query *history.HistoryFilter) (*history.HistoryStats, error) {
	return s.repo.GetStats(ctx, query)
}

// GetCountByStatus returns count of history records by status
func (s *Service) GetCountByStatus(ctx context.Context, teamID string) (map[string]int, error) {
	if teamID == "" {
		return nil, fmt.Errorf("team id is required")
	}
	return s.repo.GetCountByStatus(ctx, teamID)
}

// GetRecentActivity retrieves recent history activity
func (s *Service) GetRecentActivity(ctx context.Context, teamID string, limit int) ([]history.History, error) {
	if teamID == "" {
		return nil, fmt.Errorf("team id is required")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetRecentActivity(ctx, teamID, limit)
}
