package models

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewErrorResponse(err string, msg string, code int) ErrorResponse {
	return ErrorResponse{
		Error:   err,
		Message: msg,
		Code:    code,
	}
}

func NewSuccessResponse(msg string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: msg,
		Data:    data,
	}
}

// Common error messages
const (
	ErrInvalidRequest   = "invalid_request"
	ErrNotFound         = "not_found"
	ErrUnauthorized     = "unauthorized"
	ErrForbidden        = "forbidden"
	ErrInternalServer   = "internal_server_error"
	ErrBadRequest       = "bad_request"
)

// Common success messages
const (
	MsgCreated          = "created_successfully"
	MsgUpdated          = "updated_successfully"
	MsgDeleted          = "deleted_successfully"
	MsgPublished        = "published_successfully"
	MsgSaved            = "saved_successfully"
)