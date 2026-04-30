package generate

import (
	"fmt"
	"os"

	"seo-backend/internal/database"
	"seo-backend/internal/security"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetModelConfig(modelID, serviceType string) (*ModelConfig, error) {
	var config ModelConfig
	var encryptedKey string

	query := `
		SELECT 
			ak.key_encrypted,
			m.name,
			m.request_template,
			m.response_text_path,
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

	err := database.GetDB().QueryRow(query, modelID, serviceType).Scan(
		&encryptedKey, &config.ModelName, &config.Template,
		&config.ResponsePath, &config.BaseURL, &config.AuthType,
		&config.AuthHeader, &config.AuthPrefix, &config.Endpoint,
	)

	if err != nil {
		return nil, fmt.Errorf("model not found or no active api key: %w", err)
	}

	// Decrypt API key
	decryptor := security.NewDecryptor(os.Getenv("ENCRYPTION_KEY"))
	config.APIKey, err = decryptor.Decrypt(encryptedKey)
	if err != nil {
		// Log but continue - might be plain text
		config.APIKey = encryptedKey
	}

	return &config, nil
}
