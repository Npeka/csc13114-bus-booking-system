package middleware

import (
	"bus-booking/shared/constants"
	sharedcontext "bus-booking/shared/context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestContext(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract or generate request ID
		requestID := c.GetHeader(constants.XRequestID)
		if requestID == "" {
			requestID = sharedcontext.GenerateRequestID()
		}
		sharedcontext.SetRequestID(c, requestID)

		// Extract user information from headers (set by gateway/auth service)
		if userID := c.GetHeader(constants.XUserID); userID != "" {
			sharedcontext.SetUserID(c, userID)
		}
		if userRole := c.GetHeader(constants.XUserRole); userRole != "" {
			if roleInt, err := strconv.Atoi(userRole); err == nil {
				sharedcontext.SetUserRole(c, roleInt)
			}
		}
		if userEmail := c.GetHeader(constants.XUserEmail); userEmail != "" {
			sharedcontext.SetUserEmail(c, userEmail)
		}
		if accessToken := c.GetHeader(constants.XAccessToken); accessToken != "" {
			sharedcontext.SetAccessToken(c, accessToken)
		}

		// Set service name
		sharedcontext.SetServiceName(c, serviceName)

		// Add request context to standard context
		reqCtx := sharedcontext.GetRequestContext(c)
		ctx := sharedcontext.WithRequestContext(c.Request.Context(), reqCtx)
		c.Request = c.Request.WithContext(ctx)

		// Set response headers
		c.Header(constants.XRequestID, requestID)
		c.Header(constants.XServiceName, serviceName)

		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := sharedcontext.GetUserID(c)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": constants.ErrUnauthorized,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRole(allowedRoles ...constants.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := sharedcontext.GetUserRole(c)
		if userRole == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": constants.ErrUnauthorized,
				},
			})
			c.Abort()
			return
		}

		if !userRole.HasAnyRole(allowedRoles) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"message": constants.ErrForbidden,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
