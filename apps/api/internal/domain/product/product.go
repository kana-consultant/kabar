package product

import "time"

type UserContext interface {
	GetUserID() string
	GetTeamID() string
	GetRole() string
	IsAdmin() bool
}

// ProductFilters for filtering products
type ProductFilters struct {
	Status     string
	SyncStatus string
	Platform   string
	TeamID     string
	UserID     string
	Search     string
	OrderBy    string
	Limit      int
	Offset     int
}

type ConnectionTestResult struct {
	Success     bool   `json:"success"`
	StatusCode  int    `json:"status_code"`
	StatusText  string `json:"status_text"`
	Message     string `json:"message"`
	ProductName string `json:"product_name"`
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	Response    string `json:"response"`
	TestedAt    string `json:"tested_at"`
}

type ProductBasicInfo struct {
	ID            string
	Name          string
	CustomHeaders string
	Platform      string
	APIEndpoint   string
	APIKey        string
}

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
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Platform        string         `json:"platform"`
	APIEndpoint     string         `json:"apiEndpoint"`
	APIKeyEncrypted string         `json:"apiKey"`
	Status          ProductStatus  `json:"status"`
	LastSync        *time.Time     `json:"lastSync,omitempty"`
	SyncStatus      SyncStatus     `json:"syncStatus"`
	CreatedBy       *string        `json:"createdBy,omitempty"`
	TeamID          *string        `json:"teamId,omitempty"`
	UserID          *string        `json:"userId,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	AdapterConfig   *AdapterConfig `json:"adapterConfig,omitempty"`
}

type AdapterConfig struct {
	ID             string    `json:"id"`
	ProductID      string    `json:"productId"`
	EndpointPath   string    `json:"endpointPath"`
	HTTPMethod     string    `json:"httpMethod"`
	CustomHeaders  string    `json:"customHeaders"`
	FieldMapping   string    `json:"fieldMapping"`  // ✅ JSON template untuk konten
	MetaConfig     string    `json:"metaConfig"`    // ✅ NEW: Meta tags configuration
	SitemapConfig  string    `json:"sitemapConfig"` // ✅ NEW: Sitemap configuration
	TimeoutSeconds int       `json:"timeoutSeconds"`
	RetryCount     int       `json:"retryCount"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
type CreateProductRequest struct {
	Name          string                      `json:"name"`
	Platform      string                      `json:"platform"`
	APIEndpoint   string                      `json:"apiEndpoint"`
	APIKey        string                      `json:"apiKey"`
	TeamID        string                      `json:"teamId,omitempty"`
	AdapterConfig *CreateAdapterConfigRequest `json:"adapterConfig,omitempty"`
}

type CreateAdapterConfigRequest struct {
	EndpointPath   string `json:"endpointPath"`
	HTTPMethod     string `json:"httpMethod"`
	CustomHeaders  string `json:"customHeaders"`
	FieldMapping   string `json:"fieldMapping"`
	MetaConfig     string `json:"metaConfig"`    // ✅ TAMBAH INI
	SitemapConfig  string `json:"sitemapConfig"` // ✅ TAMBAH INI
	TimeoutSeconds int    `json:"timeoutSeconds,omitempty"`
	RetryCount     int    `json:"retryCount,omitempty"`
}
