package middleware

import (
	"bus-booking/shared/constants"
	"github.com/gin-gonic/gin"
)

// RequireRoleStringMiddleware creates role middleware from string roles (for YAML config compatibility)
func RequireRoleStringMiddleware(allowedRoleStrings ...string) gin.HandlerFunc {
	// Convert string roles to UserRole constants
	allowedRoles := constants.FromStringSlice(allowedRoleStrings)
	return RequireRoleMiddleware(allowedRoles...)
}