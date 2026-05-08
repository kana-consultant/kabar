// internal/domain/generate/repository.go
package generate

import (
	"context"
	"time"
)

type Repository interface {
	// Jika perlu menyimpan history generate ke database
	SaveGenerationHistory(ctx context.Context, history *GenerationHistory) error
	GetGenerationHistory(ctx context.Context, userID string) ([]GenerationHistory, error)
}

type GenerationHistory struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // "article" or "image"
	Topic     string    `json:"topic"`
	Prompt    string    `json:"prompt"`
	Result    string    `json:"result"`
	ModelID   string    `json:"model_id"`
	CreatedAt time.Time `json:"created_at"`
}
