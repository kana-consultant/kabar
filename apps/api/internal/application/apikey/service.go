package apikey

import (
	"context"
	"errors"
	"fmt"

	"seo-backend/internal/domain/apikey"
	"seo-backend/internal/models"
	"seo-backend/internal/pkg/crypto"
)

// Service - application layer for API Key business logic
type Service struct {
	repo      apikey.Repository
	encryptor crypto.Encryptor
}

// NewService - constructor
func NewService(repo apikey.Repository, encryptor crypto.Encryptor) apikey.Service {
	return &Service{
		repo:      repo,
		encryptor: encryptor,
	}
}

// CreateAPIKey - create new API key
func (s *Service) CreateAPIKey(ctx context.Context, req apikey.CreateAPIKeyRequest, userCtx models.UserContext) (string, error) {
	// Business validation
	if err := s.validateCreateRequest(req); err != nil {
		return "", err
	}

	// Encrypt the API key
	encryptedKey, err := s.encryptor.Encrypt(req.Key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt API key: %w", err)
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare entity
	teamID := s.getTeamIDPtr(userCtx.GetTeamID())

	key := &apikey.APIKey{
		Service:      req.Service,
		ProviderID:   req.ProviderID,
		ModelID:      req.ModelID,
		KeyEncrypted: encryptedKey,
		SystemPrompt: req.SystemPrompt,
		IsActive:     true,
		TeamID:       teamID,
		CreatedBy:    userCtx.GetUserID(),
	}

	// Create via repository
	id, err := s.repo.Create(ctx, tx, key)
	if err != nil {
		return "", fmt.Errorf("failed to create API key: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

// GetAPIKeyByID - get API key by ID
func (s *Service) GetAPIKeyByID(ctx context.Context, id string) (*apikey.APIKey, error) {
	if id == "" {
		return nil, errors.New("API key id is required")
	}

	key, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if key == nil {
		return nil, errors.New("API key not found")
	}

	// Don't return encrypted key to client
	key.KeyEncrypted = ""

	return key, nil
}

// GetAllAPIKeys - get all API keys with filters
func (s *Service) GetAllAPIKeys(ctx context.Context, userCtx models.UserContext) ([]apikey.APIKeyDetail, error) {
	keys, err := s.repo.GetAll(ctx, userCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get API keys: %w", err)
	}
	return keys, nil
}

// UpdateAPIKey - update API key
func (s *Service) UpdateAPIKey(ctx context.Context, id string, req apikey.UpdateAPIKeyRequest, userCtx models.UserContext) error {
	if id == "" {
		return errors.New("API key id is required")
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if exists
	existingKey, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}
	if existingKey == nil {
		return errors.New("API key not found")
	}

	// Business rule: check permission
	if existingKey.CreatedBy != userCtx.GetUserID() {
		teamID := s.getTeamIDPtr(userCtx.GetTeamID())
		if existingKey.TeamID == nil || teamID == nil || *existingKey.TeamID != *teamID {
			return errors.New("access denied")
		}
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Service != nil {
		updates["service"] = *req.Service
	}
	if req.ProviderID != nil {
		updates["providerId"] = *req.ProviderID
	}
	if req.ModelID != nil {
		updates["modelId"] = *req.ModelID
	}
	if req.IsActive != nil {
		updates["isActive"] = *req.IsActive
	}
	if req.SystemPrompt != nil {
		updates["systemPrompt"] = *req.SystemPrompt
	}

	// Update via repository
	if len(updates) > 0 {
		if err := s.repo.Update(ctx, tx, id, updates); err != nil {
			return fmt.Errorf("failed to update API key: %w", err)
		}
	}

	// Handle key update separately (encrypted)
	if req.Key != nil && *req.Key != "" {
		encryptedKey, err := s.encryptor.Encrypt(*req.Key)
		if err != nil {
			return fmt.Errorf("failed to encrypt API key: %w", err)
		}
		if err := s.repo.UpdateKey(ctx, tx, id, encryptedKey); err != nil {
			return fmt.Errorf("failed to update API key: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteAPIKey - delete API key
func (s *Service) DeleteAPIKey(ctx context.Context, id string, userCtx models.UserContext) error {
	if id == "" {
		return errors.New("API key id is required")
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if exists and authorize
	existingKey, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}
	if existingKey == nil {
		return errors.New("API key not found")
	}

	// Business rule: check permission
	if existingKey.CreatedBy != userCtx.GetUserID() {
		teamID := s.getTeamIDPtr(userCtx.GetTeamID())
		if existingKey.TeamID == nil || teamID == nil || *existingKey.TeamID != *teamID {
			return errors.New("access denied")
		}
	}

	// Delete via repository
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper methods
func (s *Service) validateCreateRequest(req apikey.CreateAPIKeyRequest) error {
	if req.Service == "" {
		return errors.New("service is required")
	}
	if req.ProviderID == "" {
		return errors.New("provider id is required")
	}
	if req.ModelID == "" {
		return errors.New("model id is required")
	}
	if req.Key == "" {
		return errors.New("API key is required")
	}
	return nil
}

func (s *Service) getTeamIDPtr(teamID string) *string {
	if teamID != "" && teamID != "00000000-0000-0000-0000-000000000000" {
		return &teamID
	}
	return nil
}
