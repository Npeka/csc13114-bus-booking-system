package router

import (
	"github.com/gin-gonic/gin"

	"bus-booking/shared/middleware"
	"bus-booking/template-service/internal/handler"
)

// RouterConfig holds router dependencies
type RouterConfig struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
	ServiceName string
}

// SetupRoutes configures all routes for the service
func SetupRoutes(router *gin.Engine, cfg *RouterConfig) {
	if cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router.Use(middleware.RequestContextMiddleware(cfg.ServiceName))
	router.Use(middleware.SetupCORS(nil)) // Pass nil for default CORS config
	router.Use(middleware.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": cfg.ServiceName,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", cfg.AuthHandler.Login)
			auth.POST("/refresh", cfg.AuthHandler.RefreshToken)
			auth.POST("/reset-password", cfg.AuthHandler.ResetPassword)
			auth.POST("/confirm-reset-password", cfg.AuthHandler.ConfirmResetPassword)
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
			protected.POST("/change-password", cfg.AuthHandler.ChangePassword)

			// User management routes (for users to manage themselves)
			users := protected.Group("/users")
			{
				users.GET("/:id", cfg.UserHandler.GetUser)
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
				adminUsers.POST("", cfg.UserHandler.CreateUser)
				adminUsers.GET("", cfg.UserHandler.ListUsers)
				adminUsers.GET("/:id", cfg.UserHandler.GetUser)
				adminUsers.PUT("/:id", cfg.UserHandler.UpdateUser)
				adminUsers.DELETE("/:id", cfg.UserHandler.DeleteUser)
				adminUsers.PATCH("/:id/status", cfg.UserHandler.UpdateUserStatus)
			}
		}
	}
}
