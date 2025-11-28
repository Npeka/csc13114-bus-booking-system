package jobs

import (
	"context"
	"time"

	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/service"

	"github.com/rs/zerolog/log"
)

type BookingExpirationJob struct {
	bookingService service.BookingService
	lockRepo       repository.SeatLockRepository
	interval       time.Duration
}

func NewBookingExpirationJob(bookingService service.BookingService, lockRepo repository.SeatLockRepository) *BookingExpirationJob {
	return &BookingExpirationJob{
		bookingService: bookingService,
		lockRepo:       lockRepo,
		interval:       1 * time.Minute,
	}
}

func (j *BookingExpirationJob) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	log.Info().Msg("Booking expiration job started")

	for {
		select {
		case <-ticker.C:
			j.run(ctx)
		case <-ctx.Done():
			log.Info().Msg("Booking expiration job stopped")
			return
		}
	}
}

func (j *BookingExpirationJob) run(ctx context.Context) {
	// Clean expired seat locks
	if err := j.lockRepo.CleanExpiredLocks(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to clean expired locks")
	}

	// Check for expired reservations
	// if err := j.bookingService.CheckReservationExpiry(ctx); err != nil {
	// 	log.Error().Err(err).Msg("Failed to check reservation expiry")
	// }
}
