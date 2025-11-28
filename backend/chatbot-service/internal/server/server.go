package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/chatbot-service/config"
	"bus-booking/chatbot-service/internal/handler"
	"bus-booking/chatbot-service/internal/router"
	"bus-booking/chatbot-service/internal/service"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Run() {
	handler := s.buildHandler()
	server := &http.Server{
		Addr:           s.cfg.GetServerAddr(),
		Handler:        handler,
		ReadTimeout:    s.cfg.Server.ReadTimeout,
		WriteTimeout:   s.cfg.Server.WriteTimeout,
		IdleTimeout:    s.cfg.Server.IdleTimeout,
		MaxHeaderBytes: s.cfg.Server.MaxHeaderBytes,
	}

	// Start server
	go func() {
		log.Info().
			Str("service", s.cfg.ServiceName).
			Str("address", s.cfg.GetServerAddr()).
			Msg("HTTP server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	// Chờ tín hiệu stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutdown signal received, shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server forced to shutdown")
	} else {
		log.Info().Msg("HTTP server stopped gracefully")
	}
}

func (s *Server) buildHandler() http.Handler {
	r := gin.New()

	// Initialize services
	chatbotService := service.NewChatbotService(&s.cfg.OpenAI, &s.cfg.External)

	// Initialize handlers
	chatHandler := handler.NewChatHandler(chatbotService)

	// Setup routes
	handlers := &router.Handlers{
		ChatHandler: chatHandler,
	}
	router.SetupRoutes(r, s.cfg, handlers)

	return r
}

func (s *Server) Close() {
	// Cleanup resources if needed
	log.Info().Msg("Server cleanup completed")
}
