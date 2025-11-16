package initializer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/handler"
	"bus-booking/booking-service/internal/router"
)

// InitServer initializes and starts the HTTP server
func InitServer(cfg *config.Config, bookingHandler *handler.BookingHandler) error {
	log.Println("Initializing HTTP server...")

	// Setup router
	r := router.SetupRouter(cfg, bookingHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:           cfg.GetServerAddr(),
		Handler:        r,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// Channel to listen for interrupt signal to terminate server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.GetServerAddr())
		log.Printf("Environment: %s", cfg.Server.Environment)
		log.Printf("Server URL: http://%s", cfg.GetServerAddr())

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started successfully")
	log.Println("Press Ctrl+C to shutdown server")

	// Wait for interrupt signal to gracefully shutdown server
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server shutdown complete")
	return nil
}
