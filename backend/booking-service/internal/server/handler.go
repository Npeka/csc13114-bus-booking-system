package server

import (
	"bus-booking/booking-service/internal/handler"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/router"
	"bus-booking/booking-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	repositories := repository.NewRepositories(s.db.DB)

	bookingService := service.NewBookingService(repositories)

	bookingHandler := handler.NewBookingHandler(bookingService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		BookingHandler: bookingHandler,
	})
	return engine
}
