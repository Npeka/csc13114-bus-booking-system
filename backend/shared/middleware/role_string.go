package middleware

import (
	"bus-booking/shared/constants"

	"github.com/gin-gonic/gin"
)

func RequireRoleStringMiddleware(allowedRoleStrings ...string) gin.HandlerFunc {
	allowedRoles := constants.FromStringSlice(allowedRoleStrings)
	return RequireRoleMiddleware(allowedRoles...)
}
