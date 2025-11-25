package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/handler"
)

type Handlers struct {
	AuthHandler handler.AuthHandler
	UserHandler handler.UserHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContextMiddleware(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/verify-token", ginext.WrapHandler(h.AuthHandler.VerifyToken))
			auth.POST("/firebase/auth", ginext.WrapHandler(h.AuthHandler.FirebaseAuth))
			auth.POST("/register", ginext.WrapHandler(h.AuthHandler.Register))
			auth.POST("/login", ginext.WrapHandler(h.AuthHandler.Login))
			auth.POST("/refresh-token", ginext.WrapHandler(h.AuthHandler.RefreshToken))
			auth.POST("/logout", middleware.RequireAuthMiddleware(), ginext.WrapHandler(h.AuthHandler.Logout))
		}

		users := v1.Group("/users")
		users.Use(middleware.RequireAuthMiddleware())
		{
			users.GET("/profile", ginext.WrapHandler(h.UserHandler.GetProfile))
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAuthMiddleware())
		admin.Use(middleware.RequireRoleMiddleware(constants.RoleAdmin))
		{
			adminUsers := admin.Group("/users")
			{
				adminUsers.POST("", ginext.WrapHandler(h.UserHandler.CreateUser))
				adminUsers.GET("", ginext.WrapHandler(h.UserHandler.ListUsers))
				adminUsers.GET("/:id", ginext.WrapHandler(h.UserHandler.GetUser))
				adminUsers.PUT("/:id", ginext.WrapHandler(h.UserHandler.UpdateUser))
				adminUsers.DELETE("/:id", ginext.WrapHandler(h.UserHandler.DeleteUser))
				adminUsers.PATCH("/:id/status", ginext.WrapHandler(h.UserHandler.UpdateUserStatus))
			}
		}
	}
}
