package cronjob

import (
	"context"
	"time"

	"bus-booking/trip-service/internal/service"

	"github.com/rs/zerolog/log"
)

// TripStatusCronJob auto-updates trip statuses based on time
type TripStatusCronJob struct {
	tripSvc  service.TripService
	stopChan chan struct{}
}

func NewTripStatusCronJob(tripSvc service.TripService) *TripStatusCronJob {
	return &TripStatusCronJob{
		tripSvc:  tripSvc,
		stopChan: make(chan struct{}),
	}
}

// Start begins the cronjob worker - runs every 1 minute
func (c *TripStatusCronJob) Start(ctx context.Context) {
	log.Info().Msg("Starting trip status update cronjob worker")

	// Run immediately on start
	c.processUpdates(ctx)

	// Run every 1 minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Trip status cronjob context cancelled, stopping...")
			return
		case <-c.stopChan:
			log.Info().Msg("Trip status cronjob stopped")
			return
		case <-ticker.C:
			c.processUpdates(ctx)
		}
	}
}

// Stop stops the cronjob worker
func (c *TripStatusCronJob) Stop() {
	close(c.stopChan)
}

func (c *TripStatusCronJob) processUpdates(ctx context.Context) {
	if err := c.tripSvc.ProcessTripStatusUpdates(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to process trip status updates")
	} else {
		log.Debug().Msg("Successfully processed trip status updates")
	}
}
