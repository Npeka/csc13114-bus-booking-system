package context

import (
	"bus-booking/shared/constants"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextKey defines the type for context keys to avoid collisions

// RequestContext contains request-scoped information
type RequestContext struct {
	RequestID   string
	UserID      uuid.UUID
	UserRole    constants.UserRole
	UserEmail   string
	ServiceName string
	AccessToken string
}

// GetRequestContext extracts request context from Gin context
func GetRequestContext(c *gin.Context) *RequestContext {
	return &RequestContext{
		RequestID:   GetRequestID(c),
		UserID:      GetUserID(c),
		UserRole:    GetUserRole(c),
		UserEmail:   GetUserEmail(c),
		ServiceName: GetServiceName(c),
		AccessToken: GetAccessToken(c),
	}
}

// GetRequestID gets request ID from context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(constants.XRequestID); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// SetRequestID sets request ID in context
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(constants.XRequestID, requestID)
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) uuid.UUID {
	userID, _ := uuid.Parse(c.GetString(constants.XUserID))
	return userID
}

// SetUserID sets user ID in context
func SetUserID(c *gin.Context, userID string) {
	c.Set(constants.XUserID, userID)
}

// GetUserRole gets user role from context
func GetUserRole(c *gin.Context) constants.UserRole {
	if userRole, exists := c.Get(constants.XUserRole); exists {
		if role, ok := userRole.(constants.UserRole); ok {
			return role
		}
	}
	return 1 // Default to passenger role
}

// SetUserRole sets user role in context
func SetUserRole(c *gin.Context, userRole int) {
	c.Set(constants.XUserRole, constants.UserRole(userRole))
}

// GetUserEmail gets user email from context
func GetUserEmail(c *gin.Context) string {
	if userEmail, exists := c.Get(constants.XUserEmail); exists {
		if email, ok := userEmail.(string); ok {
			return email
		}
	}
	return ""
}

// SetUserEmail sets user email in context
func SetUserEmail(c *gin.Context, userEmail string) {
	c.Set(constants.XUserEmail, userEmail)
}

// GetServiceName gets service name from context
func GetServiceName(c *gin.Context) string {
	if serviceName, exists := c.Get(constants.XServiceName); exists {
		if name, ok := serviceName.(string); ok {
			return name
		}
	}
	return ""
}

// SetServiceName sets service name in context
func SetServiceName(c *gin.Context, serviceName string) {
	c.Set(constants.XServiceName, serviceName)
}

// GenerateRequestID generates a new request ID
func GenerateRequestID() string {
	return uuid.New().String()
}

// WithRequestContext adds request context to standard context
func WithRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	ctx = context.WithValue(ctx, constants.XRequestID, reqCtx.RequestID)
	ctx = context.WithValue(ctx, constants.XUserID, reqCtx.UserID.String())
	ctx = context.WithValue(ctx, constants.XUserRole, reqCtx.UserRole)
	ctx = context.WithValue(ctx, constants.XUserEmail, reqCtx.UserEmail)
	ctx = context.WithValue(ctx, constants.XServiceName, reqCtx.ServiceName)
	return ctx
}

// FromRequestContext extracts request context from standard context
func FromRequestContext(ctx context.Context) *RequestContext {
	reqCtx := &RequestContext{}

	if requestID, ok := ctx.Value(constants.XRequestID).(string); ok {
		reqCtx.RequestID = requestID
	}

	if userIDStr, ok := ctx.Value(constants.XUserID).(string); ok && userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			reqCtx.UserID = userID
		}
	}

	if userRole, ok := ctx.Value(constants.XUserRole).(constants.UserRole); ok {
		reqCtx.UserRole = userRole
	}
	if userEmail, ok := ctx.Value(constants.XUserEmail).(string); ok {
		reqCtx.UserEmail = userEmail
	}
	if serviceName, ok := ctx.Value(constants.XServiceName).(string); ok {
		reqCtx.ServiceName = serviceName
	}

	return reqCtx
}

// GetAccessToken gets access token from context
func GetAccessToken(c *gin.Context) string {
	if accessToken, exists := c.Get("access_token"); exists {
		if token, ok := accessToken.(string); ok {
			return token
		}
	}
	return ""
}

// SetAccessToken sets access token in context
func SetAccessToken(c *gin.Context, accessToken string) {
	c.Set(constants.XAccessToken, accessToken)
}
