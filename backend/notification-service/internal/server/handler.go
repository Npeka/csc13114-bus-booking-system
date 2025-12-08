package server

import (
	"bus-booking/notification-service/internal/handler"
	"bus-booking/notification-service/internal/router"
	"bus-booking/notification-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	notificationService := service.NewNotificationService()
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
