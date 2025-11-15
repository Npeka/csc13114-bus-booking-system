package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"csc13114-bus-ticket-booking-system/shared/config"
	"csc13114-bus-ticket-booking-system/shared/db"
	"csc13114-bus-ticket-booking-system/shared/logger"
	"csc13114-bus-ticket-booking-system/shared/utils"
	"csc13114-bus-ticket-booking-system/shared/validator"
	"csc13114-bus-ticket-booking-system/template-service/internal/handler"
	"csc13114-bus-ticket-booking-system/template-service/internal/model"
	"csc13114-bus-ticket-booking-system/template-service/internal/repository"
	"csc13114-bus-ticket-booking-system/template-service/internal/router"
	"csc13114-bus-ticket-booking-system/template-service/internal/service"
)

const ServiceName = "template-service"

// Application holds the application dependencies
type Application struct {
	Config     *config.Config
	Database   *db.DatabaseManager
	Redis      *db.RedisManager
	JWTManager *utils.JWTManager

	HTTPServer *http.Server

	// Repositories
	UserRepo repository.UserRepositoryInterface

	// Services
	UserService service.UserServiceInterface
	AuthService service.AuthServiceInterface

	// Handlers
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logger
	if err := logger.SetupLogger(&cfg.Log); err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}

	log.Info().Str("service", ServiceName).Msg("Starting Template Service...")

	// Initialize validator
	validator.InitValidator()

	// Create application instance
	app := &Application{
		Config: cfg,
	}

	// Initialize dependencies
	if err := app.initDependencies(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize dependencies")
	}

	// Setup HTTP server
	app.setupHTTPServer()

	// Start server
	app.start()
}

// initDependencies initializes all application dependencies
func (app *Application) initDependencies() error {
	// Initialize database connection
	dbManager, err := db.NewPostgresConnection(&app.Config.Database, app.Config.Server.Environment)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	app.Database = dbManager

	// Initialize Redis connection
	redisManager, err := db.NewRedisConnection(&app.Config.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	app.Redis = redisManager

	// Run database migrations
	if err := app.runMigrations(); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Initialize utilities
	app.JWTManager = utils.NewJWTManager(&app.Config.JWT)

	// Initialize repositories
	app.UserRepo = repository.NewUserRepository(app.Database.DB)

	// Initialize services
	app.UserService = service.NewUserService(app.UserRepo)
	app.AuthService = service.NewAuthService(app.UserRepo, app.JWTManager)

	// Initialize handlers
	app.UserHandler = handler.NewUserHandler(app.UserService)
	app.AuthHandler = handler.NewAuthHandler(app.AuthService)

	log.Info().Msg("All dependencies initialized successfully")
	return nil
}

// runMigrations runs database migrations
func (app *Application) runMigrations() error {
	return app.Database.AutoMigrate(
		&model.User{},
	)
}

// setupHTTPServer configures the HTTP server and routes
func (app *Application) setupHTTPServer() {
	// Set Gin mode based on environment
	if app.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	ginRouter := gin.New()

	// Setup routes using the router package
	routerConfig := &router.RouterConfig{
		UserHandler: app.UserHandler,
		AuthHandler: app.AuthHandler,
		ServiceName: ServiceName,
	}
	router.SetupRoutes(ginRouter, routerConfig)

	// Create HTTP server
	app.HTTPServer = &http.Server{
		Addr:           app.Config.GetServerAddr(),
		Handler:        ginRouter,
		ReadTimeout:    app.Config.Server.ReadTimeout,
		WriteTimeout:   app.Config.Server.WriteTimeout,
		IdleTimeout:    app.Config.Server.IdleTimeout,
		MaxHeaderBytes: app.Config.Server.MaxHeaderBytes,
	}

	log.Info().Str("address", app.Config.GetServerAddr()).Msg("HTTP server configured")
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
