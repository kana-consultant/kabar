// internal/interfaces/api/common/success.go
package common

// SuccessResponse represents a success response with ID
type SuccessResponse struct {
	ID      string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message string `json:"message" example:"Operation completed successfully"`
}

// SuccessMessage represents a simple success message
type SuccessMessage struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// SuccessCreated represents a 201 created response
type SuccessCreated struct {
	ID      string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message string `json:"message" example:"Resource created successfully"`
}

// SuccessUpdated represents a successful update response
type SuccessUpdated struct {
	ID      string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message string `json:"message" example:"Resource updated successfully"`
}

// SuccessDeleted represents a successful delete response
type SuccessDeleted struct {
	ID      string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message string `json:"message" example:"Resource deleted successfully"`
}
