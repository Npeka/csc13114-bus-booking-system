package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bus-booking/booking-service/internal/client"
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/model/trip"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/shared/queue"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TripReminderJob struct {
	bookingRepo        repository.BookingRepository
	delayedQueue       queue.DelayedQueueManager
	notificationClient client.NotificationClient
	tripClient         client.TripClient
	userClient         client.UserClient
	interval           time.Duration
}

func NewTripReminderJob(
	bookingRepo repository.BookingRepository,
	delayedQueue queue.DelayedQueueManager,
	notificationClient client.NotificationClient,
	tripClient client.TripClient,
	userClient client.UserClient,
) *TripReminderJob {
	return &TripReminderJob{
		bookingRepo:        bookingRepo,
		delayedQueue:       delayedQueue,
		notificationClient: notificationClient,
		tripClient:         tripClient,
		userClient:         userClient,
		interval:           5 * time.Second,
	}
}

func (j *TripReminderJob) Start(ctx context.Context) {
	log.Info().Msg("Trip reminder job started")
	go j.runQueuePolling(ctx)
}

func (j *TripReminderJob) runQueuePolling(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
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

func (j *TripReminderJob) processQueue(ctx context.Context) {
	items, err := j.delayedQueue.Poll(ctx, "trip_reminder", 10)
	if err != nil {
		log.Error().Err(err).Msg("Failed to poll trip_reminder queue")
		return
	}

	for _, item := range items {
		j.processItem(ctx, item)
	}
}

func (j *TripReminderJob) processItem(ctx context.Context, item *queue.DelayedItem) {
	// Parse booking ID from payload
	var bookingID uuid.UUID

	payloadBytes, err := json.Marshal(item.Payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal payload for processing")
		return
	}

	if err := json.Unmarshal(payloadBytes, &bookingID); err != nil {
		log.Error().Err(err).Str("payload", string(payloadBytes)).Msg("Failed to unmarshal booking ID from payload")
		return
	}

	if err := j.sendReminder(ctx, bookingID); err != nil {
		log.Error().Err(err).Str("booking_id", bookingID.String()).Msg("Failed to send trip reminder")
	} else {
		log.Info().Str("booking_id", bookingID.String()).Msg("Successfully sent trip reminder")
	}
}

func (j *TripReminderJob) sendReminder(ctx context.Context, bookingID uuid.UUID) error {
	booking, err := j.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	user, err := j.userClient.GetUserByID(ctx, booking.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	tripData, err := j.tripClient.GetTripByID(ctx, trip.GetTripByIDRequest{
		PreLoadRoute: true,
		PreloadBus:   true,
	}, booking.TripID)
	if err != nil {
		return fmt.Errorf("failed to get trip: %w", err)
	}

	seats, err := j.tripClient.ListSeatsByIDs(ctx, getSeatIDs(booking))
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get seat details")
		// Continue potentially? Or fail? Let's continue with basic info.
	}

	seatNumbers := make([]string, len(seats))
	for i, s := range seats {
		seatNumbers[i] = s.SeatNumber
	}

	pickupPoint := "N/A"
	if tripData.Route != nil {
		pickupPoint = tripData.Route.Origin
	}

	busPlate := "Pending"
	if tripData.Bus != nil {
		busPlate = tripData.Bus.PlateNumber
	}

	req := &client.TripReminderRequest{
		Email:             user.Email,
		PassengerName:     user.FullName,
		BookingReference:  booking.BookingReference,
		DepartureLocation: getOrigin(tripData),
		Destination:       getDestination(tripData),
		DepartureTime:     tripData.DepartureTime.Format("15:04 02/01/2006"),
		SeatNumbers:       strings.Join(seatNumbers, ", "),
		BusPlate:          busPlate,
		PickupPoint:       pickupPoint,
		TicketLink:        fmt.Sprintf("http://localhost:3000/tickets/%s", booking.ID), // TODO: Use config for frontend URL
	}

	return j.notificationClient.SendTripReminder(ctx, req)
}

func getSeatIDs(booking *model.Booking) []uuid.UUID {
	ids := make([]uuid.UUID, len(booking.BookingSeats))
	for i, s := range booking.BookingSeats {
		ids[i] = s.SeatID
	}
	return ids
}

func getOrigin(trip *trip.Trip) string {
	if trip.Route != nil {
		return trip.Route.Origin
	}
	return "Unknown"
}

func getDestination(trip *trip.Trip) string {
	if trip.Route != nil {
		return trip.Route.Destination
	}
	return "Unknown"
}
