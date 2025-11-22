package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/initializer"
	"bus-booking/booking-service/internal/router"
	sharedDB "bus-booking/shared/db"
	sharedMiddleware "bus-booking/shared/middleware"
)

type Application struct {
	Config     *config.Config
	Database   *sharedDB.DatabaseManager
	Redis      *sharedDB.RedisManager
	HTTPServer *http.Server

	Services *initializer.ServiceDependencies
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	app := &Application{Config: cfg}

	// Initialize database
	app.Database, err = sharedDB.NewPostgresConnection(&cfg.BaseConfig.Database, cfg.BaseConfig.Server.Environment)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer func() {
		if err := app.Database.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	// Initialize Redis
	app.Redis, err = sharedDB.NewRedisConnection(&cfg.BaseConfig.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Redis")
	}
	defer func() {
		if err := app.Redis.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Redis connection")
		}
	}()

	// Initialize services
	app.Services = initializer.InitServices(cfg, app.Database)

	// Setup HTTP server
	app.HTTPServer = setupHTTPServer(cfg, app.Services)

	// Start server
	startServer(app)
}

func setupHTTPServer(cfg *config.Config, services *initializer.ServiceDependencies) *http.Server {
	ginRouter := router.SetupRouter(cfg, services.BookingHandler)

	// Add shared middleware
	ginRouter.Use(sharedMiddleware.SetupCORS(&cfg.BaseConfig.CORS))

	return &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      ginRouter,
		ReadTimeout:  cfg.BaseConfig.Server.ReadTimeout,
		WriteTimeout: cfg.BaseConfig.Server.WriteTimeout,
		IdleTimeout:  cfg.BaseConfig.Server.IdleTimeout,
	}
}

func startServer(app *Application) {
	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		log.Info().
			Str("service", app.Config.ServiceName).
			Str("address", app.Config.GetServerAddr()).
			Msg("Starting HTTP server...")

		if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed to start")
		}
	}()

	log.Info().
		Str("service", app.Config.ServiceName).
		Str("address", app.Config.GetServerAddr()).
		Msg("Booking Service started successfully")

	// Wait for interrupt signal
	<-quit
	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
		os.Exit(1)
	}

	log.Info().Msg("Server exited gracefully")
}
