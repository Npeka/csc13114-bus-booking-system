package cronjob

import (
	"context"
	"time"

	"bus-booking/trip-service/internal/service"

	"github.com/rs/zerolog/log"
)

// TripRescheduleCronJob auto-reschedules completed trips
type TripRescheduleCronJob struct {
	tripSvc  service.TripService
	stopChan chan struct{}
}

func NewTripRescheduleCronJob(tripSvc service.TripService) *TripRescheduleCronJob {
	return &TripRescheduleCronJob{
		tripSvc:  tripSvc,
		stopChan: make(chan struct{}),
	}
}

// Start begins the cronjob worker - runs every 30 minutes
func (c *TripRescheduleCronJob) Start(ctx context.Context) {
	log.Info().Msg("Starting trip reschedule cronjob worker")

	// Run immediately on start
	c.rescheduleCompletedTrips(ctx)

	// Run every 30 minutes
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Trip reschedule cronjob context cancelled, stopping...")
			return
		case <-c.stopChan:
			log.Info().Msg("Trip reschedule cronjob stopped")
			return
		case <-ticker.C:
			c.rescheduleCompletedTrips(ctx)
		}
	}
}

// Stop stops the cronjob worker
func (c *TripRescheduleCronJob) Stop() {
	close(c.stopChan)
}

// rescheduleCompletedTrips finds completed trips and reschedules them
func (c *TripRescheduleCronJob) rescheduleCompletedTrips(ctx context.Context) {
	log.Info().Msg("Checking for completed trips to reschedule")

	// Find trips that completed in the last 24 hours
	completedTrips, err := c.tripSvc.GetCompletedTripsForReschedule(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get completed trips for reschedule")
		return
	}

	if len(completedTrips) == 0 {
		log.Debug().Msg("No completed trips to reschedule")
		return
	}

	successCount := 0
	failureCount := 0

	for _, trip := range completedTrips {
		// Calculate the time gap between departure and arrival (duration)
		duration := trip.ArrivalTime.Sub(trip.DepartureTime)

		// Find the next week's same day/time
		nextDeparture := trip.DepartureTime.AddDate(0, 0, 7) // +7 days (1 week)
		nextArrival := nextDeparture.Add(duration)

		// Skip if next departure is in the past
		if nextDeparture.Before(time.Now()) {
			log.Warn().
				Str("trip_id", trip.ID.String()).
				Time("next_departure", nextDeparture).
				Msg("Skipping reschedule - next departure in past")
			continue
		}

		// Update the trip with new datetime (reschedule in-place)
		err := c.tripSvc.RescheduleTrip(ctx, trip.ID, nextDeparture, nextArrival)
		if err != nil {
			log.Error().
				Err(err).
				Str("trip_id", trip.ID.String()).
				Time("old_departure", trip.DepartureTime).
				Time("new_departure", nextDeparture).
				Msg("Failed to reschedule trip")
			failureCount++
			continue
		}

		log.Info().
			Str("trip_id", trip.ID.String()).
			Time("old_departure", trip.DepartureTime).
			Time("new_departure", nextDeparture).
			Dur("duration", duration).
			Msg("Successfully rescheduled trip")
		successCount++
	}

	log.Info().
		Int("total", len(completedTrips)).
		Int("success", successCount).
		Int("failed", failureCount).
		Msg("Completed trip rescheduling")
}
