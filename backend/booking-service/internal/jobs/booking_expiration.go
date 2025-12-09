package jobs

import (
	"context"
	"encoding/json"
	"time"

	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/queue"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BookingExpirationJob struct {
	bookingService service.BookingService
	lockRepo       repository.SeatLockRepository
	delayedQueue   queue.DelayedQueueManager
	interval       time.Duration
}

func NewBookingExpirationJob(
	bookingService service.BookingService,
	lockRepo repository.SeatLockRepository,
	delayedQueue queue.DelayedQueueManager,
) *BookingExpirationJob {
	return &BookingExpirationJob{
		bookingService: bookingService,
		lockRepo:       lockRepo,
		delayedQueue:   delayedQueue,
		interval:       1 * time.Minute,
	}
}

func (j *BookingExpirationJob) Start(ctx context.Context) {
	log.Info().Msg("Booking expiration job started")

	// 1. Start seat lock cleanup ticker
	go j.runLockCleanup(ctx)

	// 2. Start queue polling
	go j.runQueuePolling(ctx)
}

func (j *BookingExpirationJob) runLockCleanup(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := j.lockRepo.CleanExpiredLocks(ctx); err != nil {
				log.Error().Err(err).Msg("Failed to clean expired locks")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (j *BookingExpirationJob) runQueuePolling(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // Poll every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.processQueue(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (j *BookingExpirationJob) processQueue(ctx context.Context) {
	items, err := j.delayedQueue.Poll(ctx, "booking_expiry", 10) // Limit 10 items per poll
	if err != nil {
		log.Error().Err(err).Msg("Failed to poll booking_expiry queue")
		return
	}

	for _, item := range items {
		j.processItem(ctx, item)
	}
}

func (j *BookingExpirationJob) processItem(ctx context.Context, item *queue.DelayedItem) {
	// Parse booking ID from payload
	// Payload was saved as booking.ID (uuid.UUID) which marshals to string
	var bookingID uuid.UUID

	// Handle if payload is string or object. uuid marshals to string.
	// We need to marshal/unmarshal or type assert carefully.
	// Since we used json.Marshal in Schedule, item.Payload is interface{}.
	// If it was unmarshaled from JSON in Poll, it depends on how encoding/json treated it.
	// UUIDs are strings in JSON.

	payloadBytes, err := json.Marshal(item.Payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal payload for processing")
		return
	}

	if err := json.Unmarshal(payloadBytes, &bookingID); err != nil {
		log.Error().Err(err).Str("payload", string(payloadBytes)).Msg("Failed to unmarshal booking ID from payload")
		return
	}

	if err := j.bookingService.ExpireBooking(ctx, bookingID); err != nil {
		log.Error().Err(err).Str("booking_id", bookingID.String()).Msg("Failed to expire booking")
	} else {
		log.Info().Str("booking_id", bookingID.String()).Msg("Successfully processed booking expiration")
	}
}
