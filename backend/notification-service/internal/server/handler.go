package server

import (
	"bus-booking/notification-service/internal/handler"
	"bus-booking/notification-service/internal/router"
	"bus-booking/notification-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Server) buildHandler() http.Handler {
	emailService, err := service.NewEmailService(s.cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create email service")
	}

	notificationService := service.NewNotificationService(emailService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		NotificationHandler: notificationHandler,
	})
	return engine
}
