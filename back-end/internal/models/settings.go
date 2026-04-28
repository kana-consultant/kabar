package models

import (
	"encoding/json"
	"time"
)

type Setting struct {
	ID          string          `json:"id"`
	Key         string          `json:"key"`
	Value       json.RawMessage `json:"value"`
	Description *string         `json:"description,omitempty"`
	UpdatedBy   *string         `json:"updatedBy,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type UpdateSettingRequest struct {
	Value       json.RawMessage `json:"value"`
	Description *string         `json:"description,omitempty"`
}

// Default settings values
type DefaultPostStatus struct {
	Value string `json:"value"`
}

type AutoSaveDraft struct {
	Enabled bool `json:"enabled"`
}

type TimezoneSetting struct {
	Timezone string `json:"timezone"`
}

type EmailNotifications struct {
	Enabled bool `json:"enabled"`
}