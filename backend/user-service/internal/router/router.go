package router

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/handler"
	"bus-booking/user-service/internal/repository"
)

type RouterConfig struct {
	Config       *config.Config
	FirebaseAuth *auth.Client
	UserHandler  handler.UserHandler
	AuthHandler  handler.AuthHandler
	UserRepo     repository.UserRepository
}

func SetupRoutes(router *gin.Engine, cfg *RouterConfig) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.Config.CORS))
	router.Use(middleware.RequestContextMiddleware(cfg.Config.ServiceName))
	router.GET(health.Path, health.Handler(cfg.Config.ServiceName))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/verify-token", ginext.WrapHandler(cfg.AuthHandler.VerifyToken))
			auth.POST("/firebase/auth", ginext.WrapHandler(cfg.AuthHandler.FirebaseAuth))
			auth.POST("/refresh-token", middleware.RequireAuthMiddleware(), ginext.WrapHandler(cfg.AuthHandler.RefreshToken))
			auth.POST("/logout", middleware.RequireAuthMiddleware(), ginext.WrapHandler(cfg.AuthHandler.Logout))
		}

		users := v1.Group("/users")
		users.Use(middleware.RequireAuthMiddleware())
		{
			users.GET("/profile", ginext.WrapHandler(cfg.UserHandler.GetProfile))
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAuthMiddleware())
		admin.Use(middleware.RequireRoleMiddleware(constants.RoleAdmin))
		{
			adminUsers := admin.Group("/users")
			{
				adminUsers.POST("", ginext.WrapHandler(cfg.UserHandler.CreateUser))
				adminUsers.GET("", ginext.WrapHandler(cfg.UserHandler.ListUsers))
				adminUsers.GET("/:id", ginext.WrapHandler(cfg.UserHandler.GetUser))
				adminUsers.PUT("/:id", ginext.WrapHandler(cfg.UserHandler.UpdateUser))
				adminUsers.DELETE("/:id", ginext.WrapHandler(cfg.UserHandler.DeleteUser))
				adminUsers.PATCH("/:id/status", ginext.WrapHandler(cfg.UserHandler.UpdateUserStatus))
			}
		}
	}
}
