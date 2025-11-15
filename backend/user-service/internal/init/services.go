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

// InitServices initializes all application services and handlers
func InitServices(cfg *config.Config, database *sharedDB.DatabaseManager, firebaseAuth *auth.Client) (*handler.UserHandler, *handler.AuthHandler) {
	// Initialize utilities
	jwtManager := utils.NewJWTManager(&cfg.JWT)

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtManager, firebaseAuth, cfg)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	log.Info().Msg("All services and handlers initialized successfully")
	return userHandler, authHandler
}
