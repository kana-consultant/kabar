// internal/models/api_provider.go
package models

import (
	"database/sql"
	"time"
)

type APIProvider struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	DisplayName       string         `json:"displayName"`
	Description       string         `json:"description"`
	BaseURL           string         `json:"baseUrl"`
	AuthType          string         `json:"authType"`
	AuthHeader        string         `json:"authHeader"`
	AuthPrefix        string         `json:"authPrefix"`
	TextEndpoint      string         `json:"textEndpoint"`
	ImageEndpoint     sql.NullString `json:"imageEndpoint"`
	DefaultHeaders    sql.NullString `json:"defaultHeaders"`
	RequestTemplate   sql.NullString `json:"requestTemplate"`
	ResponseTextPath  string         `json:"responseTextPath"`
	ResponseImagePath sql.NullString `json:"responseImagePath"`
	IsActive          bool           `json:"isActive"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}
