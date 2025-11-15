package response

import (
	"net/http"
	"strconv"

	sharedcontext "bus-booking/shared/context"

	"github.com/gin-gonic/gin"
)

const (
	// HeaderRequestID represents the X-Request-ID header key
	HeaderRequestID = "X-Request-ID"
)

// Response represents the standardized API response format
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
}

// Error represents error details in the response
type Error struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Stack   string            `json:"stack,omitempty"`
}

// Meta represents metadata for pagination and other info
type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Total      int64       `json:"total,omitempty"`
	Version    string      `json:"version,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := Response{
		Data: data,
	}

	// Add X-Request-ID header
	c.Header(HeaderRequestID, getRequestID(c))
	c.JSON(http.StatusOK, response)
}

// CreatedResponse sends a created response
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	response := Response{
		Data: data,
	}

	// Add X-Request-ID header
	c.Header(HeaderRequestID, getRequestID(c))
	c.JSON(http.StatusCreated, response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	response := Response{
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}

	// Add X-Request-ID header
	c.Header(HeaderRequestID, getRequestID(c))
	c.JSON(statusCode, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors []ValidationError) {
	details := make(map[string]string)
	for _, err := range errors {
		details[err.Field] = err.Message
	}

	response := Response{
		Error: &Error{
			Code:    "VALIDATION_ERROR",
			Message: "One or more validation errors occurred",
			Details: details,
		},
	}

	// Add X-Request-ID header
	c.Header(HeaderRequestID, getRequestID(c))
	c.JSON(http.StatusBadRequest, response)
}

// BadRequestResponse sends a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// ForbiddenResponse sends a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Access denied"
	}
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", message)
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

// TooManyRequestsResponse sends a too many requests response
func TooManyRequestsResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Too many requests"
	}
	ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", message)
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *gin.Context, message string, data interface{}, pagination *Pagination) {
	response := Response{
		Data: data,
		Meta: &Meta{
			Pagination: pagination,
			Total:      pagination.Total,
		},
	}

	// Add X-Request-ID header
	c.Header("X-Request-ID", getRequestID(c))
	c.JSON(http.StatusOK, response)
}

// NoContentResponse sends a no content response
func NoContentResponse(c *gin.Context) {
	// Add X-Request-ID header
	c.Header(HeaderRequestID, getRequestID(c))
	c.Status(http.StatusNoContent)
}

// getRequestID gets the request ID from context
func getRequestID(c *gin.Context) string {
	return sharedcontext.GetRequestID(c)
}

// CalculatePagination calculates pagination information
func CalculatePagination(page, limit int, total int64) *Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// GetPaginationFromQuery extracts pagination parameters from query
func GetPaginationFromQuery(c *gin.Context) (page, limit int) {
	page = 1
	limit = 10

	if p := c.DefaultQuery("page", "1"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.DefaultQuery("limit", "10"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit
}
