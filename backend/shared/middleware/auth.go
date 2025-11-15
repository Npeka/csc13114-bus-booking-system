package middleware

import (
	"strings"
	"time"

	"bus-booking/shared/config"
	"bus-booking/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuthMiddleware creates a JWT authentication middleware
func JWTAuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.UnauthorizedResponse(c, "Missing authorization header")
			c.Abort()
			return
		}

		// Check if token has Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.SecretKey), nil
		})

		if err != nil {
			log.Debug().Err(err).Msg("JWT parsing error")
			response.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		// Check if token is valid
		if !token.Valid {
			response.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			response.UnauthorizedResponse(c, "Invalid token claims")
			c.Abort()
			return
		}

		// Check token expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			response.UnauthorizedResponse(c, "Token expired")
			c.Abort()
			return
		}

		// Check issuer and audience
		if claims.Issuer != cfg.Issuer {
			response.UnauthorizedResponse(c, "Invalid token issuer")
			c.Abort()
			return
		}

		if len(claims.Audience) > 0 && claims.Audience[0] != cfg.Audience {
			response.UnauthorizedResponse(c, "Invalid token audience")
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// OptionalJWTAuthMiddleware creates an optional JWT authentication middleware
func OptionalJWTAuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Check if token has Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Invalid header format, continue without authentication
			c.Next()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(*Claims); ok {
			// Check token expiration
			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
				c.Next()
				return
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)
			c.Set("jwt_claims", claims)
		}

		c.Next()
	}
}

// RequireRole creates a role-based authorization middleware
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.UnauthorizedResponse(c, "Authentication required")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			response.ForbiddenResponse(c, "Invalid user role")
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		roleAllowed := false
		for _, role := range allowedRoles {
			if roleStr == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			log.Warn().
				Str("user_role", roleStr).
				Strs("allowed_roles", allowedRoles).
				Str("user_id", c.GetString("user_id")).
				Msg("Access denied due to insufficient role")

			response.ForbiddenResponse(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOwnership creates an ownership-based authorization middleware
func RequireOwnership(resourceOwnerFunc func(*gin.Context) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			response.UnauthorizedResponse(c, "Authentication required")
			c.Abort()
			return
		}

		// Get resource owner ID
		ownerID, err := resourceOwnerFunc(c)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get resource owner")
			response.InternalServerErrorResponse(c, "Failed to verify ownership")
			c.Abort()
			return
		}

		// Check if user owns the resource or has admin role
		userRole := c.GetString("user_role")
		if userID != ownerID && userRole != "admin" && userRole != "super_admin" {
			log.Warn().
				Str("user_id", userID).
				Str("owner_id", ownerID).
				Str("role", userRole).
				Msg("Access denied due to ownership check")

			response.ForbiddenResponse(c, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext gets user ID from Gin context
func GetUserIDFromContext(c *gin.Context) string {
	return c.GetString("user_id")
}

// GetUserEmailFromContext gets user email from Gin context
func GetUserEmailFromContext(c *gin.Context) string {
	return c.GetString("user_email")
}

// GetUserRoleFromContext gets user role from Gin context
func GetUserRoleFromContext(c *gin.Context) string {
	return c.GetString("user_role")
}

// GetClaimsFromContext gets JWT claims from Gin context
func GetClaimsFromContext(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*Claims)
	return jwtClaims, ok
}
