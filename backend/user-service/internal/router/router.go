package router

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/shared/middleware"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/handler"
	"bus-booking/user-service/internal/repository"
)

type RouterConfig struct {
	ServiceName  string
	Config       *config.Config
	FirebaseAuth *auth.Client
	UserHandler  handler.UserHandler
	AuthHandler  handler.AuthHandler
	UserRepo     repository.UserRepository
}

func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	router.Use(middleware.RequestContextMiddleware(config.ServiceName))
	router.Use(middleware.SetupCORS(&config.Config.CORS))
	router.Use(middleware.Logger())

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
			auth.POST("/verify-token", ginext.WrapHandler(config.AuthHandler.VerifyToken))
			auth.POST("/firebase/auth", ginext.WrapHandler(config.AuthHandler.FirebaseAuth))
			auth.POST("/refresh-token", middleware.RequireAuthMiddleware(), ginext.WrapHandler(config.AuthHandler.RefreshToken))
			auth.POST("/logout", middleware.RequireAuthMiddleware(), ginext.WrapHandler(config.AuthHandler.Logout))
		}

		users := v1.Group("/users")
		users.Use(middleware.RequireAuthMiddleware())
		{
			users.GET("/profile", ginext.WrapHandler(config.UserHandler.GetProfile))
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAuthMiddleware())
		admin.Use(middleware.RequireRoleMiddleware(constants.RoleAdmin))
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
