package service

import (
	"context"
	"time"

	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
)

type SeatService interface {
	GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error)
	ReserveSeat(ctx context.Context, req *model.ReserveSeatRequest) error
	ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error
	CheckReservationExpiry(ctx context.Context) error
}

type SeatServiceImpl struct {
	seatStatusRepo repository.SeatStatusRepository
}

func NewSeatService(seatStatusRepo repository.SeatStatusRepository) SeatService {
	return &SeatServiceImpl{
		seatStatusRepo: seatStatusRepo,
	}
}

// GetSeatAvailability retrieves seat availability for a trip
func (s *SeatServiceImpl) GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error) {
	seatStatuses, err := s.seatStatusRepo.GetSeatStatusByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	var availableSeats, reservedSeats, bookedSeats []uuid.UUID
	seatMap := make(map[uuid.UUID]*model.SeatStatus)

	for _, seatStatus := range seatStatuses {
		seatMap[seatStatus.SeatID] = seatStatus

		switch seatStatus.Status {
		case "available":
			availableSeats = append(availableSeats, seatStatus.SeatID)
		case "reserved":
			// Check if reservation has expired
			if seatStatus.ReservedUntil != nil && time.Now().UTC().After(*seatStatus.ReservedUntil) {
				// Release expired reservation
				if err := s.seatStatusRepo.ReleaseSeat(ctx, tripID, seatStatus.SeatID); err != nil {
					log.Error().Err(err).Msg("Failed to release expired seat reservation")
				} else {
					availableSeats = append(availableSeats, seatStatus.SeatID)
				}
			} else {
				reservedSeats = append(reservedSeats, seatStatus.SeatID)
			}
		case "booked":
			bookedSeats = append(bookedSeats, seatStatus.SeatID)
		}
	}

	return &model.SeatAvailabilityResponse{
		TripID:         tripID,
		AvailableSeats: availableSeats,
		ReservedSeats:  reservedSeats,
		BookedSeats:    bookedSeats,
		SeatDetails:    seatMap,
	}, nil
}

// ReserveSeat reserves a seat for a user
func (s *SeatServiceImpl) ReserveSeat(ctx context.Context, req *model.ReserveSeatRequest) error {
	// Check if seat is available
	seatStatuses, err := s.seatStatusRepo.GetSeatStatusByTripID(ctx, req.TripID)
	if err != nil {
		return err
	}

	for _, seatStatus := range seatStatuses {
		if seatStatus.SeatID == req.SeatID {
			if seatStatus.Status != "available" {
				return ginext.NewConflictError("seat is not available")
			}
			break
		}
	}

	// Default reservation time is 15 minutes
	reservationTime := 15 * time.Minute
	if req.ReservationMinutes > 0 {
		reservationTime = time.Duration(req.ReservationMinutes) * time.Minute
	}

	return s.seatStatusRepo.ReserveSeat(ctx, req.TripID, req.SeatID, req.UserID, reservationTime)
}

// ReleaseSeat releases a reserved seat
func (s *SeatServiceImpl) ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error {
	return s.seatStatusRepo.ReleaseSeat(ctx, tripID, seatID)
}

// CheckReservationExpiry checks and releases expired seat reservations
func (s *SeatServiceImpl) CheckReservationExpiry(ctx context.Context) error {
	// This should be called periodically by a background job
	// For now, it's a placeholder implementation
	log.Info().Msg("Checking for expired seat reservations")

	// In a real implementation, you would:
	// 1. Query all reserved seats with expired reservations
	// 2. Release them back to available status
	// 3. Optionally notify users about expired reservations

	return nil
}
