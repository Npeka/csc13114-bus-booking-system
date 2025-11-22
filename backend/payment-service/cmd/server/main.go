package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog/log"

	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/initializer"
	sharedDB "bus-booking/shared/db"
	"bus-booking/shared/validator"
)

type Application struct {
	Config       *config.Config
	Database     *sharedDB.DatabaseManager
	Redis        *sharedDB.RedisManager
	FirebaseAuth *auth.Client
	HTTPServer   *http.Server

	Services *initializer.ServiceDependencies
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	if err := initializer.SetupLogger(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}

	log.Info().Str("service", cfg.ServiceName).Msg("Starting User Service...")

	validator.InitValidator()

	app := &Application{Config: cfg}
	if err := app.initDependencies(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize dependencies")
	}

	app.setupHTTPServer()
	app.start()
}

func (app *Application) initDependencies() error {
	// Initialize database
	var err error
	app.Database, err = initializer.InitDatabase(app.Config)
	if err != nil {
		return err
	}

	// Initialize services and handlers
	app.Services = initializer.InitServices(app.Config, app.Database)

	log.Info().Msg("All dependencies initialized successfully")
	return nil
}

// setupHTTPServer configures the HTTP server and routes
func (app *Application) setupHTTPServer() {
	app.HTTPServer = initializer.InitHTTPServer(app.Config, app.Services)
}

// start starts the HTTP server with graceful shutdown
func (app *Application) start() {
	// Start server in a goroutine
	go func() {
		log.Info().Str("address", app.Config.GetServerAddr()).Msg("Starting HTTP server...")
		if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	// Close database connections
	if app.Database != nil {
		if err := app.Database.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}

	// Close Redis connection
	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Redis connection")
		}
	}

	log.Info().Msg("Server shutdown complete")
}
