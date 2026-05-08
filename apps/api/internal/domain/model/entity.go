// internal/domain/aimodel/entity.go
package aimodel

type AIModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	ProviderID  string `json:"providerId"`
	IsActive    bool   `json:"isActive"`
}

type ModelWithStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProviderID  string `json:"providerId"`
	DisplayName string `json:"displayName"`
}
