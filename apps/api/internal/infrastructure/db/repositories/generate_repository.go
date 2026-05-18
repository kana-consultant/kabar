// internal/infrastructure/repository/generate/repository.go
package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"seo-backend/internal/domain/generate"
	"seo-backend/internal/security"
)

type GenerateRepositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) generate.Repository {
	return &GenerateRepositoryImpl{db: db}
}

func (r *GenerateRepositoryImpl) GetModelConfig(ctx context.Context, modelID, serviceType string) (*generate.ModelConfig, error) {
	var config generate.ModelConfig
	var encryptedKey string

	query := `
	SELECT 
		ak.key_encrypted,
		ak.system_prompt,
		m.name,
		m.request_template,
		COALESCE(m.response_text_path, '') AS response_text_path,
		COALESCE(m.response_image_path, '') AS response_image_path,
		p.base_url,
		p.auth_type,
		p.auth_header,
		p.auth_prefix,
		p.text_endpoint
	FROM api_keys ak
	JOIN ai_models m ON ak.model_id = m.id
	JOIN api_providers p ON m.provider_id = p.id
	WHERE ak.id = $1 AND ak.is_active = true AND ak.service = $2
	LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, modelID, serviceType).Scan(
		&encryptedKey,
		&config.SystemPrompt,
		&config.ModelName,
		&config.Template,
		&config.ResponsePath,
		&config.ResponseImagePath,
		&config.BaseURL,
		&config.AuthType,
		&config.AuthHeader,
		&config.AuthPrefix,
		&config.Endpoint,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("model not found or no active api key for model_id: %s, service: %s", modelID, serviceType)
		}
		return nil, fmt.Errorf("failed to get model config: %w", err)
	}

	// Decrypt API key
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable is not set")
	}

	decryptor := security.NewDecryptor(encryptionKey)
	decryptedKey, err := decryptor.Decrypt(encryptedKey)
	if err != nil {
		// Log warning but don't fail - maybe the key is not encrypted
		// Use the encrypted key as is if decryption fails
		config.APIKey = encryptedKey
	} else {
		config.APIKey = decryptedKey
	}

	return &config, nil
}
func (r *GenerateRepositoryImpl) SaveHistory(ctx context.Context, history *generate.GenerationHistory) error {
	query := `
		INSERT INTO generation_history (type, topic, result, model_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		history.Type, history.Topic, history.Result, history.ModelID, history.CreatedAt,
	)
	return err
}
