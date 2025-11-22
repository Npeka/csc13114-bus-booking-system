package initializer

import (
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/router"
)

func InitHTTPServer(
	cfg *config.Config,
	firebaseAuth *auth.Client,
	services *ServiceDependencies,
) *http.Server {
	if cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	ginRouter := gin.New()

	routerConfig := &router.RouterConfig{
		Config:       cfg,
		FirebaseAuth: firebaseAuth,
		UserHandler:  services.UserHandler,
		AuthHandler:  services.AuthHandler,
		UserRepo:     services.UserRepo,
	}
	router.SetupRoutes(ginRouter, routerConfig)

	server := &http.Server{
		Addr:           cfg.GetServerAddr(),
		Handler:        ginRouter,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	log.Info().Str("address", cfg.GetServerAddr()).Msg("HTTP server configured")
	return server
}
