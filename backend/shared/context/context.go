package context

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextKey defines the type for context keys to avoid collisions
type ContextKey string

const (
	// Request context keys
	RequestIDKey   ContextKey = "request_id"
	UserIDKey      ContextKey = "user_id"
	UserRoleKey    ContextKey = "user_role"
	UserEmailKey   ContextKey = "user_email"
	ServiceNameKey ContextKey = "service_name"
)

// RequestContext contains request-scoped information
type RequestContext struct {
	RequestID   string
	UserID      string
	UserRole    string
	UserEmail   string
	ServiceName string
}

// GetRequestContext extracts request context from Gin context
func GetRequestContext(c *gin.Context) *RequestContext {
	return &RequestContext{
		RequestID:   GetRequestID(c),
		UserID:      GetUserID(c),
		UserRole:    GetUserRole(c),
		UserEmail:   GetUserEmail(c),
		ServiceName: GetServiceName(c),
	}
}

// GetRequestID gets request ID from context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(string(RequestIDKey)); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// SetRequestID sets request ID in context
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(string(RequestIDKey), requestID)
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get(string(UserIDKey)); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// SetUserID sets user ID in context
func SetUserID(c *gin.Context, userID string) {
	c.Set(string(UserIDKey), userID)
}

// GetUserRole gets user role from context
func GetUserRole(c *gin.Context) string {
	if userRole, exists := c.Get(string(UserRoleKey)); exists {
		if role, ok := userRole.(string); ok {
			return role
		}
	}
	return ""
}

// SetUserRole sets user role in context
func SetUserRole(c *gin.Context, userRole string) {
	c.Set(string(UserRoleKey), userRole)
}

// GetUserEmail gets user email from context
func GetUserEmail(c *gin.Context) string {
	if userEmail, exists := c.Get(string(UserEmailKey)); exists {
		if email, ok := userEmail.(string); ok {
			return email
		}
	}
	return ""
}

// SetUserEmail sets user email in context
func SetUserEmail(c *gin.Context, userEmail string) {
	c.Set(string(UserEmailKey), userEmail)
}

// GetServiceName gets service name from context
func GetServiceName(c *gin.Context) string {
	if serviceName, exists := c.Get(string(ServiceNameKey)); exists {
		if name, ok := serviceName.(string); ok {
			return name
		}
	}
	return ""
}

// SetServiceName sets service name in context
func SetServiceName(c *gin.Context, serviceName string) {
	c.Set(string(ServiceNameKey), serviceName)
}

// GenerateRequestID generates a new request ID
func GenerateRequestID() string {
	return uuid.New().String()
}

// WithRequestContext adds request context to standard context
func WithRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	ctx = context.WithValue(ctx, RequestIDKey, reqCtx.RequestID)
	ctx = context.WithValue(ctx, UserIDKey, reqCtx.UserID)
	ctx = context.WithValue(ctx, UserRoleKey, reqCtx.UserRole)
	ctx = context.WithValue(ctx, UserEmailKey, reqCtx.UserEmail)
	ctx = context.WithValue(ctx, ServiceNameKey, reqCtx.ServiceName)
	return ctx
}

// FromRequestContext extracts request context from standard context
func FromRequestContext(ctx context.Context) *RequestContext {
	reqCtx := &RequestContext{}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		reqCtx.RequestID = requestID
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		reqCtx.UserID = userID
	}
	if userRole, ok := ctx.Value(UserRoleKey).(string); ok {
		reqCtx.UserRole = userRole
	}
	if userEmail, ok := ctx.Value(UserEmailKey).(string); ok {
		reqCtx.UserEmail = userEmail
	}
	if serviceName, ok := ctx.Value(ServiceNameKey).(string); ok {
		reqCtx.ServiceName = serviceName
	}

	return reqCtx
}
