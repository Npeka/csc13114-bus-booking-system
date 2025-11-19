package appinit

import (
	"firebase.google.com/go/v4/auth"

	sharedDB "bus-booking/shared/db"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/handler"
	"bus-booking/user-service/internal/repository"
	"bus-booking/user-service/internal/service"
	"bus-booking/user-service/internal/utils"

	"github.com/rs/zerolog/log"
)

type ServiceDependencies struct {
	UserHandler handler.UserHandler
	AuthHandler handler.AuthHandler
	UserRepo    repository.UserRepository
}

func InitServices(cfg *config.Config, database *sharedDB.DatabaseManager, redis *sharedDB.RedisManager, firebaseAuth *auth.Client) *ServiceDependencies {
	jwtManager := utils.NewJWTManager(&cfg.JWT)

	userRepo := repository.NewUserRepository(database.DB)
	tokenBlacklistMgr := service.NewTokenBlacklistManager(redis, jwtManager)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtManager, firebaseAuth, cfg, tokenBlacklistMgr)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	log.Info().Msg("All services and handlers initialized successfully")

	return &ServiceDependencies{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		UserRepo:    userRepo,
	}
}
