package appinit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/handler"
	"bus-booking/user-service/internal/router"
)

func InitHTTPServer(cfg *config.Config, userHandler *handler.UserHandler, authHandler *handler.AuthHandler, serviceName string) *http.Server {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	ginRouter := gin.New()

	routerConfig := &router.RouterConfig{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		ServiceName: serviceName,
		Config:      cfg,
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
