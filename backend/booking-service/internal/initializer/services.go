package initializer

import (
	"log"

	"gorm.io/gorm"

	"bus-booking/booking-service/internal/handler"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/service"
)

// InitServices initializes all services and dependencies
func InitServices(db *gorm.DB) (*handler.BookingHandler, error) {
	log.Println("Initializing services...")

	// Initialize repository
	bookingRepo := repository.NewBookingRepository(db)

	// Initialize service
	bookingService := service.NewBookingService(bookingRepo)

	// Initialize handler
	bookingHandler := handler.NewBookingHandler(bookingService)

	log.Println("Services initialized successfully")
	return bookingHandler, nil
}
