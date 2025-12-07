package service

import (
	"context"
	"time"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatStatusService interface {
	InitSeatsForTrip(ctx context.Context, tripID uuid.UUID, seats []model.SeatInitData) error
	GetSeatAvailability(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error)
	CheckAndReleaseSeat(ctx context.Context, tripID uuid.UUID, seatID uuid.UUID) error
}

type SeatStatusServiceImpl struct {
	seatStatusRepo repository.SeatStatusRepository
}

func NewSeatStatusService(seatStatusRepo repository.SeatStatusRepository) SeatStatusService {
	return &SeatStatusServiceImpl{
		seatStatusRepo: seatStatusRepo,
	}
}

// InitSeatsForTrip initializes seat statuses when a trip is created
func (s *SeatStatusServiceImpl) InitSeatsForTrip(ctx context.Context, tripID uuid.UUID, seats []model.SeatInitData) error {
	// Check if seats already initialized
	existing, err := s.seatStatusRepo.GetSeatStatusByTripID(ctx, tripID)
	if err != nil {
		return ginext.NewInternalServerError("failed to check existing seats")
	}

	if len(existing) > 0 {
		log.Warn().Str("trip_id", tripID.String()).Msg("seats already initialized for this trip")
		return ginext.NewBadRequestError("seats already initialized for this trip")
	}

	// Create seat status records
	var seatStatuses []*model.SeatStatus
	for _, seat := range seats {
		seatStatus := &model.SeatStatus{
			TripID:     tripID,
			SeatID:     seat.SeatID,
			SeatNumber: seat.SeatNumber,
			Status:     "available",
		}
		seatStatuses = append(seatStatuses, seatStatus)
	}

	// Bulk insert
	if err := s.seatStatusRepo.BulkUpdateSeatStatus(ctx, seatStatuses); err != nil {
		return ginext.NewInternalServerError("failed to initialize seats")
	}

	log.Info().
		Str("trip_id", tripID.String()).
		Int("seat_count", len(seatStatuses)).
		Msg("seats initialized successfully")

	return nil
}

// GetSeatAvailability returns all seats for a trip with their current status
func (s *SeatStatusServiceImpl) GetSeatAvailability(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error) {
	seats, err := s.seatStatusRepo.GetSeatStatusByTripID(ctx, tripID)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get seat availability")
	}

	// Check for expired reservations and release them
	for _, seat := range seats {
		if seat.Status == "reserved" && seat.ReservedUntil != nil {
			if time.Now().UTC().After(*seat.ReservedUntil) {
				// Release expired reservation
				if err := s.seatStatusRepo.ReleaseSeat(ctx, tripID, seat.SeatID); err != nil {
					log.Error().Err(err).Str("seat_id", seat.SeatID.String()).Msg("failed to release expired seat")
				} else {
					seat.Status = "available"
					seat.ReservedUntil = nil
					seat.UserID = nil
				}
			}
		}
	}

	return seats, nil
}

// CheckAndReleaseSeat releases a seat if its reservation has expired
func (s *SeatStatusServiceImpl) CheckAndReleaseSeat(ctx context.Context, tripID uuid.UUID, seatID uuid.UUID) error {
	seats, err := s.seatStatusRepo.GetSeatStatusByTripID(ctx, tripID)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		if seat.SeatID == seatID && seat.Status == "reserved" && seat.ReservedUntil != nil {
			if time.Now().UTC().After(*seat.ReservedUntil) {
				return s.seatStatusRepo.ReleaseSeat(ctx, tripID, seatID)
			}
		}
	}

	return nil
}
