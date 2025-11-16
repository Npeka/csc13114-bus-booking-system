package appinit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/router"
)

func InitHTTPServer(
	cfg *config.Config,
	serviceName string,
	services *ServiceDependencies,
) *http.Server {
	if cfg.BaseConfig.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	ginRouter := gin.New()

	routerConfig := &router.RouterConfig{
		Config:       cfg,
		ServiceName:  serviceName,
		TripHandler:  services.TripHandler,
		RouteHandler: services.RouteHandler,
		BusHandler:   services.BusHandler,
	}
	router.SetupRoutes(ginRouter, routerConfig)

	server := &http.Server{
		Addr:           cfg.GetServerAddr(),
		Handler:        ginRouter,
		ReadTimeout:    cfg.BaseConfig.Server.ReadTimeout,
		WriteTimeout:   cfg.BaseConfig.Server.WriteTimeout,
		IdleTimeout:    cfg.BaseConfig.Server.IdleTimeout,
		MaxHeaderBytes: cfg.BaseConfig.Server.MaxHeaderBytes,
	}

	log.Info().Str("address", cfg.GetServerAddr()).Msg("HTTP server configured")
	return server
}
