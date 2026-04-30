package apikey

import (
	"database/sql"
	"time"
)

type APIKeyDetail struct {
	ID                  string    `json:"id"`
	Service             string    `json:"service"`
	ProviderID          string    `json:"providerId"`
	ModelID             string    `json:"modelId"`
	IsActive            bool      `json:"isActive"`
	SystemPrompt        string    `json:"systemPrompt"`
	CreatedBy           *string   `json:"createdBy"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	ProviderName        string    `json:"providerName"`
	ProviderDisplayName string    `json:"providerDisplayName"`
	ModelName           string    `json:"modelName"`
	ModelDisplayName    string    `json:"modelDisplayName"`
}

type APIKey struct {
	ID           string    `json:"id"`
	Service      string    `json:"service"`
	ProviderID   string    `json:"providerId"`
	ModelID      string    `json:"modelId"`
	IsActive     bool      `json:"isActive"`
	SystemPrompt string    `json:"systemPrompt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// helper scan nullable
func ScanString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
