package constants

// Error messages
const (
	ErrInvalidRequestData = "Invalid request data"
	ErrUnauthorized       = "Unauthorized access"
	ErrForbidden          = "Access forbidden"
	ErrNotFound           = "Resource not found"
	ErrInternalServer     = "Internal server error"
	ErrBadRequest         = "Bad request"
	ErrConflict           = "Resource conflict"
	ErrValidationFailed   = "Validation failed"
	ErrRateLimitExceeded  = "Rate limit exceeded"
	ErrRequestTimeout     = "Request timeout"
	ErrServiceUnavailable = "Service unavailable"
)

// Success messages
const (
	MsgOperationSuccess = "Operation completed successfully"
	MsgCreatedSuccess   = "Resource created successfully"
	MsgUpdatedSuccess   = "Resource updated successfully"
	MsgDeletedSuccess   = "Resource deleted successfully"
)

// Error codes
const (
	CodeInvalidRequest     = "INVALID_REQUEST"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "NOT_FOUND"
	CodeInternalError      = "INTERNAL_ERROR"
	CodeBadRequest         = "BAD_REQUEST"
	CodeConflict           = "CONFLICT"
	CodeValidationFailed   = "VALIDATION_FAILED"
	CodeRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
	CodeRequestTimeout     = "REQUEST_TIMEOUT"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)
