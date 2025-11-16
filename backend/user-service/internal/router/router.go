package router

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	"bus-booking/shared/ginext"
	"bus-booking/shared/middleware"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/handler"
	localMiddleware "bus-booking/user-service/internal/middleware"
	"bus-booking/user-service/internal/repository"
)

type RouterConfig struct {
	UserHandler  *handler.UserHandler
	AuthHandler  *handler.AuthHandler
	ServiceName  string
	Config       *config.Config
	FirebaseAuth *auth.Client
	UserRepo     repository.UserRepositoryInterface
}

func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	router.Use(middleware.RequestContextMiddleware(config.ServiceName))
	router.Use(middleware.SetupCORS(&config.Config.CORS))
	router.Use(middleware.Logger())

	var firebaseMiddleware *localMiddleware.FirebaseAuthMiddleware
	if config.FirebaseAuth != nil {
		firebaseMiddleware = localMiddleware.NewFirebaseAuthMiddleware(config.FirebaseAuth, config.UserRepo)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": config.ServiceName,
		})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", ginext.WrapHandler(config.AuthHandler.Signup))
			auth.POST("/signin", ginext.WrapHandler(config.AuthHandler.Signin))
			auth.POST("/oauth2/signin", ginext.WrapHandler(config.AuthHandler.OAuth2Signin))
			auth.POST("/refresh", ginext.WrapHandler(config.AuthHandler.RefreshToken))
			auth.POST("/signout", ginext.WrapHandler(config.AuthHandler.Signout))
			auth.POST("/verify-token", ginext.WrapHandler(config.AuthHandler.VerifyToken))
		}

		if firebaseMiddleware != nil {
			firebase := v1.Group("/firebase")
			firebase.Use(firebaseMiddleware.FirebaseAuth())
			{
				firebase.GET("/profile", func(c *gin.Context) {
					user, exists := c.Get("user")
					if !exists {
						c.JSON(500, gin.H{"error": "User not found in context"})
						return
					}
					c.JSON(200, gin.H{"data": user, "message": "Profile retrieved successfully"})
				})
			}
		}

		protected := v1.Group("")
		protected.Use(middleware.RequireAuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Profile endpoint - TODO: implement"})
			})
			protected.PUT("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Update profile endpoint - TODO: implement"})
			})

			users := protected.Group("/users")
			{
				users.GET("/:id", ginext.WrapHandler(config.UserHandler.GetUser))
			}
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAuthMiddleware())
		admin.Use(middleware.RequireRoleMiddleware("admin"))
		{
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
