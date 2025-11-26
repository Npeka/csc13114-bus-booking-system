package server

import (
	"bus-booking/user-service/internal/handler"
	"bus-booking/user-service/internal/repository"
	"bus-booking/user-service/internal/router"
	"bus-booking/user-service/internal/service"
	"bus-booking/user-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	jwtManager := utils.NewJWTManager(&s.cfg.JWT)

	userRepo := repository.NewUserRepository(s.db.DB)

	tokenBlacklistMgr := service.NewTokenBlacklistManager(s.rd, jwtManager)
	passwordResetService := service.NewPasswordResetService(s.rd)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(s.cfg, jwtManager, s.fa, tokenBlacklistMgr, userRepo, passwordResetService)

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
