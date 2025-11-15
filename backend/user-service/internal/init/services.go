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

func InitServices(cfg *config.Config, database *sharedDB.DatabaseManager, firebaseAuth *auth.Client) (*handler.UserHandler, *handler.AuthHandler) {
	jwtManager := utils.NewJWTManager(&cfg.JWT)

	userRepo := repository.NewUserRepository(database.DB)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtManager, firebaseAuth, cfg)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	log.Info().Msg("All services and handlers initialized successfully")
	return userHandler, authHandler
}
