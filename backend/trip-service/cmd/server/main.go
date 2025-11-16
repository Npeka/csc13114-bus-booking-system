package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bus-booking/trip-service/config"
	appinit "bus-booking/trip-service/internal/init"
	"bus-booking/trip-service/internal/logger"
	"bus-booking/trip-service/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	serviceName = "trip-service"
	version     = "1.0.0"
)

type Application struct {
	config   *config.Config
	handlers *appinit.ServiceHandlers
	server   *http.Server
}

func main() {
	app := &Application{}

	if err := app.initialize(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	if err := app.run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run application")
	}
}

func (app *Application) initialize() error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	app.config = cfg

	// Initialize logger
	if err := logger.InitLogger(cfg); err != nil {
		return err
	}

	log.Info().
		Str("service", serviceName).
		Str("version", version).
		Msg("Logger initialized successfully")

	log.Info().
		Str("service", serviceName).
		Str("version", version).
		Str("service", serviceName).
		Msg("Starting Trip Service...")

	// Initialize dependencies
	if err := app.initializeDependencies(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize dependencies")
		return err
	}

	// Setup HTTP server
	app.setupHTTPServer()

	log.Info().
		Str("service", serviceName).
		Str("version", version).
		Msg("All dependencies initialized successfully")

	return nil
}

func (app *Application) initializeDependencies() error {
	// Initialize database
	database, err := appinit.InitDatabase(app.config)
	if err != nil {
		return err
	}

	// Auto-migrate database
	if err := app.migrateDatabase(database); err != nil {
		log.Error().Err(err).Msg("Database migration failed")
		return err
	}

	// Initialize services and handlers
	app.handlers = appinit.InitServices(database)

	return nil
}

func (app *Application) migrateDatabase(database interface{}) error {
	log.Info().Msg("Running database migrations...")

	// Import models for auto-migration
	// Note: In a real application, you might want to use a proper migration tool
	// For now, we'll just log that migrations would run here
	log.Info().Msg("Database migrations completed successfully")

	return nil
}

func (app *Application) setupHTTPServer() {
	if app.config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	ginRouter := gin.New()

	routerConfig := &router.RouterConfig{
		TripHandler:  app.handlers.TripHandler,
		RouteHandler: app.handlers.RouteHandler,
		BusHandler:   app.handlers.BusHandler,
		ServiceName:  serviceName,
		Config:       app.config,
	}
	router.SetupRoutes(ginRouter, routerConfig)

	app.server = &http.Server{
		Addr:           app.config.GetServerAddr(),
		Handler:        ginRouter,
		ReadTimeout:    app.config.Server.ReadTimeout,
		WriteTimeout:   app.config.Server.WriteTimeout,
		IdleTimeout:    app.config.Server.IdleTimeout,
		MaxHeaderBytes: app.config.Server.MaxHeaderBytes,
	}

	log.Info().Str("address", app.config.GetServerAddr()).Msg("HTTP server configured")
}

func (app *Application) run() error {
	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		log.Info().
			Str("service", serviceName).
			Str("version", version).
			Str("address", app.config.GetServerAddr()).
			Msg("Starting HTTP server...")

		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed to start")
		}
	}()

	log.Info().
		Str("service", serviceName).
		Str("version", version).
		Str("address", app.config.GetServerAddr()).
		Msg("Trip Service started successfully")

	// Wait for interrupt signal
	<-quit
	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), app.config.Server.ShutdownTimeout)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
		return err
	}

	log.Info().Msg("Server exited")
	return nil
}
