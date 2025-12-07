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
	router.Use(middleware.RequestContext(cfg.ServiceName))
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
			auth.POST("/forgot-password", ginext.WrapHandler(h.AuthHandler.ForgotPassword))
			auth.POST("/reset-password", ginext.WrapHandler(h.AuthHandler.ResetPassword))
			auth.POST("/guest", ginext.WrapHandler(h.AuthHandler.CreateGuestAccount))
			auth.POST("/logout", middleware.RequireAuth(), ginext.WrapHandler(h.AuthHandler.Logout))
		}

		users := v1.Group("/users")
		users.Use(middleware.RequireAuth())
		{
			users.GET("/profile", ginext.WrapHandler(h.UserHandler.GetProfile))
			users.PUT("/profile", ginext.WrapHandler(h.UserHandler.UpdateProfile))

			users.Use(middleware.RequireRole(constants.RoleAdmin))
			{
				users.GET("", ginext.WrapHandler(h.UserHandler.ListUsers))
				users.GET("/:id", ginext.WrapHandler(h.UserHandler.GetUser))
				users.POST("", ginext.WrapHandler(h.UserHandler.CreateUser))
				users.PUT("/:id", ginext.WrapHandler(h.UserHandler.UpdateUser))
				users.DELETE("/:id", ginext.WrapHandler(h.UserHandler.DeleteUser))
			}
		}
	}
}
