package helper

// Constants for HTTP methods and defaults
const (
	DefaultHTTPMethod = "POST"
	DefaultTimeout    = 60
	DefaultRetryCount = 3
	DefaultPriority   = "0.7"
	DefaultChangefreq = "weekly"
)

// Sensitive fields for logging redaction
var sensitiveFields = map[string]bool{
	"image_url":     true,
	"image":         true,
	"photo":         true,
	"picture":       true,
	"attachment":    true,
	"file":          true,
	"body":          true,
	"description":   true,
	"password":      true,
	"token":         true,
	"secret":        true,
	"api_key":       true,
	"apikey":        true,
	"authorization": true,
	"auth":          true,
}

// Valid HTTP methods
var validHTTPMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}
