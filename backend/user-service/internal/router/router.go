package router

import (
	"github.com/gin-gonic/gin"

	"bus-booking/shared/ginext"
	"bus-booking/shared/middleware"
	"bus-booking/user-service/internal/handler"
)

// RouterConfig holds router dependencies
type RouterConfig struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
	ServiceName string
}

// SetupRoutes configures all routes for the service
func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	// Apply global middleware
	router.Use(middleware.RequestContextMiddleware(config.ServiceName))
	router.Use(middleware.SetupCORS(nil)) // Pass nil for default CORS config
	router.Use(middleware.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": config.ServiceName,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", ginext.WrapHandler(config.AuthHandler.Signup))
			auth.POST("/signin", ginext.WrapHandler(config.AuthHandler.Signin))
			auth.POST("/oauth2/signin", ginext.WrapHandler(config.AuthHandler.OAuth2Signin))
			auth.POST("/refresh", ginext.WrapHandler(config.AuthHandler.RefreshToken))
			auth.POST("/signout", ginext.WrapHandler(config.AuthHandler.Signout))
			auth.POST("/verify-token", ginext.WrapHandler(config.AuthHandler.VerifyToken))
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.RequireAuthMiddleware())
		{
			// User profile routes
			protected.GET("/profile", func(c *gin.Context) {
				// Get current user's profile
				// Implementation would get user ID from context and call user service
				c.JSON(200, gin.H{"message": "Profile endpoint - TODO: implement"})
			})
			protected.PUT("/profile", func(c *gin.Context) {
				// Update current user's profile
				c.JSON(200, gin.H{"message": "Update profile endpoint - TODO: implement"})
			})
			// Note: Change password would be implemented in user service if needed

			// User management routes (for users to manage themselves)
			users := protected.Group("/users")
			{
				users.GET("/:id", ginext.WrapHandler(config.UserHandler.GetUser))
			}
		}

		// Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAuthMiddleware())
		admin.Use(middleware.RequireRoleMiddleware("admin"))
		{
			// Admin user management
			adminUsers := admin.Group("/users")
			{
				adminUsers.POST("", ginext.WrapHandler(config.UserHandler.CreateUser))
				adminUsers.GET("", ginext.WrapHandler(config.UserHandler.ListUsers))
				adminUsers.GET("/:id", ginext.WrapHandler(config.UserHandler.GetUser))
				adminUsers.PUT("/:id", ginext.WrapHandler(config.UserHandler.UpdateUser))
				adminUsers.DELETE("/:id", ginext.WrapHandler(config.UserHandler.DeleteUser))
				adminUsers.PATCH("/:id/status", ginext.WrapHandler(config.UserHandler.UpdateUserStatus))
			}
		}
	}
}
