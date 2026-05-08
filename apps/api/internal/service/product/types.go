package product

type UserContext interface {
	GetUserID() string
	GetTeamID() string
	GetRole() string
	IsAdmin() bool
}

type ProductFilters struct {
	Platform   string
	Status     string
	SyncStatus string
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
	ID          string
	Name        string
	Platform    string
	APIEndpoint string
	APIKey      string
}

type AdapterConfig struct {
	EndpointPath   string
	HTTPMethod     string
	CustomHeaders  map[string]string
	TimeoutSeconds int
	RetryCount     int
}
