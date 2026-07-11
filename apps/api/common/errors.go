// internal/interfaces/api/common/errors.go
package common

// ErrorResponse400 untuk Bad Request
type ErrorResponse400 struct {
	Error  string `json:"error" example:"Invalid request body"`
	Status int    `json:"status" example:"400"`
}

// ErrorResponse404 untuk Not Found
type ErrorResponse404 struct {
	Error  string `json:"error" example:"Resource not found"`
	Status int    `json:"status" example:"404"`
}

// ErrorResponse409 untuk Conflict
type ErrorResponse409 struct {
	Error  string `json:"error" example:"Resource already exists"`
	Status int    `json:"status" example:"409"`
}

// ErrorResponse500 untuk Internal Server Error
type ErrorResponse500 struct {
	Error  string `json:"error" example:"Internal server error"`
	Status int    `json:"status" example:"500"`
}
