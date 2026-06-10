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

func (r *GenerateRepositoryImpl) GetModelConfig(
	ctx context.Context,
	apiKeyID string,
	serviceType string,
) (*generate.ModelConfig, error) {

	var config generate.ModelConfig
	var encryptedKey string

	query := `
	SELECT
		ak.key_encrypted,
		COALESCE(ak.system_prompt, '') AS ak_system_prompt,
		mf.name AS model_name,
		COALESCE(mf.system_prompt, '') AS mf_system_prompt,
		COALESCE(mf.max_tokens, 1024) AS max_tokens,
		COALESCE(mf.temperature, 1.0) AS temperature,
		COALESCE(rs.request_template, '') AS request_template,
		COALESCE(rs.response_text_path, '') AS response_text_path,
		COALESCE(rs.response_image_path, '') AS response_image_path,
		p.base_url,
		p.auth_type,
		p.auth_header,
		p.auth_prefix,
		COALESCE(rs.endpoint_path, '') AS endpoint_path
	FROM api_keys ak
	INNER JOIN api_providers p
		ON p.id = ak.provider_id
	INNER JOIN model_families mf
		ON mf.id = ak.model_id
	LEFT JOIN request_schemas rs
		ON rs.id = mf.schema_id
	WHERE ak.id = $1
	LIMIT 1
`

	err := r.db.QueryRowContext(
		ctx,
		query,
		apiKeyID,
		serviceType,
	).Scan(
		&encryptedKey,
		&config.APISystemPrompt,
		&config.ModelName,
		&config.ModelSystemPrompt,
		&config.MaxTokens,
		&config.Temperature,
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
			return nil, fmt.Errorf(
				"api key not found or inactive: %s, service: %s",
				apiKeyID,
				serviceType,
			)
		}
		return nil, fmt.Errorf("failed to get model config: %w", err)
	}

	// decrypt key
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable is not set")
	}

	decryptor := security.NewDecryptor(encryptionKey)

	decryptedKey, err := decryptor.Decrypt(encryptedKey)
	if err != nil {
		config.APIKey = encryptedKey
	} else {
		config.APIKey = decryptedKey
	}

	// Prioritaskan system prompt dari API key, fallback ke model family
	if config.APISystemPrompt != "" {
		config.SystemPrompt = config.APISystemPrompt
	} else {
		config.SystemPrompt = config.ModelSystemPrompt
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
