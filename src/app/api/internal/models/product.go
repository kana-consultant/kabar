package models

import (
	"time"
)

type ProductPlatform string

const (
	PlatformWordpress ProductPlatform = "wordpress"
	PlatformShopify   ProductPlatform = "shopify"
	PlatformCustom    ProductPlatform = "custom"
)

type ProductStatus string

const (
	ProductStatusConnected    ProductStatus = "connected"
	ProductStatusPending      ProductStatus = "pending"
	ProductStatusError        ProductStatus = "error"
	ProductStatusDisconnected ProductStatus = "disconnected"
)

type SyncStatus string

const (
	SyncStatusIdle    SyncStatus = "idle"
	SyncStatusSyncing SyncStatus = "syncing"
	SyncStatusSuccess SyncStatus = "success"
	SyncStatusFailed  SyncStatus = "failed"
)

type Product struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Platform        ProductPlatform `json:"platform"`
	APIEndpoint     string          `json:"apiEndpoint"`
	APIKeyEncrypted string          `json:"apiKey"`
	Status          ProductStatus   `json:"status"`
	LastSync        *time.Time      `json:"lastSync,omitempty"`
	SyncStatus      SyncStatus      `json:"syncStatus"`
	CreatedBy       *string         `json:"createdBy,omitempty"`
	TeamID          *string         `json:"teamId,omitempty"`
	UserID          *string         `json:"userId,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	AdapterConfig   *AdapterConfig  `json:"adapterConfig,omitempty"`
}

type AdapterConfig struct {
	ID             string            `json:"id"`
	ProductID      string            `json:"productId"`
	EndpointPath   string            `json:"endpointPath"`
	HTTPMethod     string            `json:"httpMethod"`
	CustomHeaders  map[string]string `json:"customHeaders"`
	FieldMapping   string            `json:"fieldMapping"`
	TimeoutSeconds int               `json:"timeoutSeconds"`
	RetryCount     int               `json:"retryCount"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type CreateProductRequest struct {
	Name          string                      `json:"name"`
	Platform      ProductPlatform             `json:"platform"`
	APIEndpoint   string                      `json:"apiEndpoint"`
	APIKey        string                      `json:"apiKey"`
	TeamID        *string                     `json:"teamId,omitempty"`
	AdapterConfig *CreateAdapterConfigRequest `json:"adapterConfig,omitempty"`
}

type CreateAdapterConfigRequest struct {
	EndpointPath   string            `json:"endpointPath"`
	HTTPMethod     string            `json:"httpMethod"`
	CustomHeaders  map[string]string `json:"customHeaders"`
	FieldMapping   string            `json:"fieldMapping"`
	TimeoutSeconds int               `json:"timeoutSeconds,omitempty"`
	RetryCount     int               `json:"retryCount,omitempty"`
}
