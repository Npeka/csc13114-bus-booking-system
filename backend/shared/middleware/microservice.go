package middleware

import (
	"bus-booking/shared/constants"
	sharedcontext "bus-booking/shared/context"

	"github.com/gin-gonic/gin"
)

// RequestContextMiddleware extracts microservice headers and sets them in context
func RequestContextMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract or generate request ID
		requestID := c.GetHeader(constants.HeaderRequestID)
		if requestID == "" {
			requestID = sharedcontext.GenerateRequestID()
		}
		sharedcontext.SetRequestID(c, requestID)

		// Extract user information from headers (set by gateway/auth service)
		if userID := c.GetHeader(constants.HeaderUserID); userID != "" {
			sharedcontext.SetUserID(c, userID)
		}
		if userRole := c.GetHeader(constants.HeaderUserRole); userRole != "" {
			sharedcontext.SetUserRole(c, userRole)
		}
		if userEmail := c.GetHeader(constants.HeaderUserEmail); userEmail != "" {
			sharedcontext.SetUserEmail(c, userEmail)
		}

		// Set service name
		sharedcontext.SetServiceName(c, serviceName)

		// Add request context to standard context
		reqCtx := sharedcontext.GetRequestContext(c)
		ctx := sharedcontext.WithRequestContext(c.Request.Context(), reqCtx)
		c.Request = c.Request.WithContext(ctx)

		// Set response headers
		c.Header(constants.HeaderRequestID, requestID)
		c.Header(constants.HeaderServiceName, serviceName)

		c.Next()
	}
}

// RequireAuthMiddleware ensures that user authentication headers are present
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := sharedcontext.GetUserID(c)
		if userID == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error": gin.H{
					"code":    constants.CodeUnauthorized,
					"message": constants.ErrUnauthorized,
				},
				"request_id": sharedcontext.GetRequestID(c),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoleMiddleware ensures that user has required role
func RequireRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := sharedcontext.GetUserRole(c)
		if userRole == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error": gin.H{
					"code":    constants.CodeUnauthorized,
					"message": constants.ErrUnauthorized,
				},
				"request_id": sharedcontext.GetRequestID(c),
			})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		roleAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(403, gin.H{
				"success": false,
				"error": gin.H{
					"code":    constants.CodeForbidden,
					"message": constants.ErrForbidden,
				},
				"request_id": sharedcontext.GetRequestID(c),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
