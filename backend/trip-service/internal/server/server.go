package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"bus-booking/shared/db"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/cronjob"
)

type Server struct {
	cfg        *config.Config
	db         *db.DatabaseManager
	redis      db.RedisManager
	cronjob    *cronjob.TripRescheduleCronJob
	statusCron *cronjob.TripStatusCronJob
}

func NewServer(
	cfg *config.Config,
	db *db.DatabaseManager,
	redis db.RedisManager,
) *Server {
	return &Server{cfg: cfg, db: db, redis: redis}
}

func (s *Server) Run() {
	handler, cronJob, statusCron := s.buildHandler()
	s.cronjob = cronJob
	s.statusCron = statusCron

	server := &http.Server{
		Addr:           s.cfg.GetServerAddr(),
		Handler:        handler,
		ReadTimeout:    s.cfg.Server.ReadTimeout,
		WriteTimeout:   s.cfg.Server.WriteTimeout,
		IdleTimeout:    s.cfg.Server.IdleTimeout,
		MaxHeaderBytes: s.cfg.Server.MaxHeaderBytes,
	}

	// Start cronjob in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Info().Msg("Starting trip schedule cronjob")
		s.cronjob.Start(ctx)
	}()

	go func() {
		log.Info().Msg("Starting trip status cronjob")
		s.statusCron.Start(ctx)
	}()

	// Start HTTP server
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

	log.Info().Msg("Shutdown signal received, shutting down services...")

	// Stop cronjob first
	log.Info().Msg("Stopping cronjob...")
	s.cronjob.Stop()
	s.statusCron.Stop()
	cancel()

	// Then stop HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("HTTP server forced to shutdown")
	} else {
		log.Info().Msg("HTTP server stopped gracefully")
	}
}

func (s *Server) Close() {
	if err := s.db.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connection")
	}
}
