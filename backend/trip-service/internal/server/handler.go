package server

import (
	"bus-booking/shared/storage"
	"bus-booking/trip-service/internal/client"
	"bus-booking/trip-service/internal/cronjob"
	"bus-booking/trip-service/internal/handler"
	"bus-booking/trip-service/internal/repository"
	"bus-booking/trip-service/internal/router"
	"bus-booking/trip-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Server) buildHandler() (http.Handler, *cronjob.TripRescheduleCronJob, *cronjob.TripStatusCronJob) {
	bookingClient := client.NewBookingClient(s.cfg.ServiceName, s.cfg.External.BookingServiceURL)
	paymentClient := client.NewPaymentClient(s.cfg.ServiceName, s.cfg.External.PaymentServiceURL)

	// Initialize repositories
	tripRepo := repository.NewTripRepository(s.db.DB)
	routeRepo := repository.NewRouteRepository(s.db.DB)
	routeStopRepo := repository.NewRouteStopRepository(s.db.DB)
	busRepo := repository.NewBusRepository(s.db.DB)
	seatRepo := repository.NewSeatRepository(s.db.DB)

	// Initialize storage service
	storageService, err := storage.NewS3StorageService(storage.S3Config{
		AccessKey: s.cfg.Storage.AccessKeyID,
		SecretKey: s.cfg.Storage.SecretAccessKey,
		Endpoint:  s.cfg.Storage.Endpoint,
		Bucket:    s.cfg.Storage.BucketName,
		Region:    s.cfg.Storage.Region,
		UseSSL:    s.cfg.Storage.UseSSL,
		CDNDomain: s.cfg.Storage.CDNDomain,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize storage service")
	}

	// Initialize services
	tripService := service.NewTripService(tripRepo, routeRepo, routeStopRepo, busRepo, seatRepo, bookingClient, paymentClient)
	routeService := service.NewRouteService(routeRepo)
	busService := service.NewBusService(busRepo, seatRepo, storageService)
	routeStopService := service.NewRouteStopService(routeStopRepo, routeRepo)
	seatService := service.NewSeatService(seatRepo)
	constantsService := service.NewConstantsService()

	// Initialize trip reschedule cronjob
	cronJob := cronjob.NewTripRescheduleCronJob(tripService)
	statusCron := cronjob.NewTripStatusCronJob(tripService)

	// Initialize handlers
	tripHandler := handler.NewTripHandler(tripService)
	routeHandler := handler.NewRouteHandler(routeService)
	busHandler := handler.NewBusHandler(busService)
	routeStopHandler := handler.NewRouteStopHandler(routeStopService)
	seatHandler := handler.NewSeatHandler(seatService)
	constantsHandler := handler.NewConstantsHandler(constantsService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		TripHandler:      tripHandler,
		RouteHandler:     routeHandler,
		BusHandler:       busHandler,
		RouteStopHandler: routeStopHandler,
		SeatHandler:      seatHandler,
		ConstantsHandler: constantsHandler,
	})
	return engine, cronJob, statusCron
}
