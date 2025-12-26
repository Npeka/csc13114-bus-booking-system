package server

import (
	"bus-booking/user-service/internal/client"
	"bus-booking/user-service/internal/handler"
	"bus-booking/user-service/internal/repository"
	"bus-booking/user-service/internal/router"
	"bus-booking/user-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	notificationClient := client.NewNotificationClient("notification-service", s.cfg.External.NotificationServiceURL)

	userRepo := repository.NewUserRepository(s.db.DB)

	jwtManager := service.NewJWTManager(&s.cfg.JWT)
	tokenManager := service.NewTokenManager(s.redis, jwtManager)
	firebaseAuth := service.NewFirebaseAuth(s.firebaseAuth)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(s.cfg, jwtManager, firebaseAuth, tokenManager, userRepo, s.redis, notificationClient)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	})
	return engine
}
