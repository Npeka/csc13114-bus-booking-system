package service

import (
	"context"
	"fmt"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatService interface {
	GetListByIDs(ctx context.Context, seatIDs []uuid.UUID) ([]model.Seat, error)
	Update(ctx context.Context, req *model.UpdateSeatRequest, id uuid.UUID) (*model.Seat, error)
}

type SeatServiceImpl struct {
	seatRepo repository.SeatRepository
}

func NewSeatService(seatRepo repository.SeatRepository) SeatService {
	return &SeatServiceImpl{
		seatRepo: seatRepo,
	}
}

func (s *SeatServiceImpl) GetListByIDs(ctx context.Context, seatIDs []uuid.UUID) ([]model.Seat, error) {
	seats, err := s.seatRepo.GetListByIDs(ctx, seatIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list seats by IDs")
		return nil, fmt.Errorf("failed to list seats: %w", err)
	}
	return seats, nil
}

func (s *SeatServiceImpl) Update(ctx context.Context, req *model.UpdateSeatRequest, id uuid.UUID) (*model.Seat, error) {
	seat, err := s.seatRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("seat_id", id.String()).Msg("Seat not found")
		return nil, fmt.Errorf("seat not found: %w", err)
	}

	// Update fields if provided
	if req.SeatNumber != nil {
		seat.SeatNumber = *req.SeatNumber
	}
	if req.Row != nil {
		seat.Row = *req.Row
	}
	if req.Column != nil {
		seat.Column = *req.Column
	}
	if req.SeatType != nil {
		seat.SeatType = *req.SeatType
		// Update price multiplier if seat type changed
		if req.PriceMultiplier == nil {
			seat.PriceMultiplier = req.SeatType.GetPriceMultiplier()
		}
	}
	if req.PriceMultiplier != nil {
		seat.PriceMultiplier = *req.PriceMultiplier
	}
	if req.IsAvailable != nil {
		seat.IsAvailable = *req.IsAvailable
	}
	if req.Floor != nil {
		seat.Floor = *req.Floor
	}

	if err := s.seatRepo.Update(ctx, seat); err != nil {
		log.Error().Err(err).Msg("Failed to update seat")
		return nil, fmt.Errorf("failed to update seat: %w", err)
	}

	log.Info().Str("seat_id", seat.ID.String()).Msg("Seat updated successfully")
	return seat, nil
}
