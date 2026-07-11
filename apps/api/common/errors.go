// internal/interfaces/api/common/errors.go
package common

// ErrorResponse400 untuk Bad Request
type ErrorResponse400 struct {
	Error  string `json:"error" example:"Invalid request body"`
	Status int    `json:"status" example:"400"`
}

// ErrorResponse401 untuk Unauthorized
type ErrorResponse401 struct {
	Error  string `json:"error" example:"Unauthorized access"`
	Status int    `json:"status" example:"401"`
}

// ErrorResponse403 untuk Forbidden
type ErrorResponse403 struct {
	Error  string `json:"error" example:"Access denied - admin role required"`
	Status int    `json:"status" example:"403"`
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

// ErrorResponse422 untuk Unprocessable Entity
type ErrorResponse422 struct {
	Error  string `json:"error" example:"Validation failed"`
	Status int    `json:"status" example:"422"`
}

// ErrorResponse429 untuk Too Many Requests
type ErrorResponse429 struct {
	Error  string `json:"error" example:"Too many requests, please try again later"`
	Status int    `json:"status" example:"429"`
}

// ErrorResponse500 untuk Internal Server Error
type ErrorResponse500 struct {
	Error  string `json:"error" example:"Internal server error"`
	Status int    `json:"status" example:"500"`
}

// ErrorResponse502 untuk Bad Gateway
type ErrorResponse502 struct {
	Error  string `json:"error" example:"Bad gateway"`
	Status int    `json:"status" example:"502"`
}

// ErrorResponse503 untuk Service Unavailable
type ErrorResponse503 struct {
	Error  string `json:"error" example:"Service temporarily unavailable"`
	Status int    `json:"status" example:"503"`
}
