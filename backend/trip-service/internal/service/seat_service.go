package service

import (
	"context"
	"fmt"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatService interface {
	ListByIDs(ctx context.Context, seatIDs []uuid.UUID) ([]model.Seat, error)
	CreateSeatsFromTemplate(ctx context.Context, req *model.BulkCreateSeatsRequest) ([]model.Seat, error)
	GetSeatMap(ctx context.Context, busID uuid.UUID) (*model.SeatMapResponse, error)
	CreateSeat(ctx context.Context, req *model.CreateSeatRequest) (*model.Seat, error)
	UpdateSeat(ctx context.Context, id uuid.UUID, req *model.UpdateSeatRequest) (*model.Seat, error)
	DeleteSeat(ctx context.Context, id uuid.UUID) error
}

type SeatServiceImpl struct {
	seatRepo repository.SeatRepository
	busRepo  repository.BusRepository
}

func NewSeatService(seatRepo repository.SeatRepository, busRepo repository.BusRepository) SeatService {
	return &SeatServiceImpl{
		seatRepo: seatRepo,
		busRepo:  busRepo,
	}
}

func (s *SeatServiceImpl) ListByIDs(ctx context.Context, seatIDs []uuid.UUID) ([]model.Seat, error) {
	seats, err := s.seatRepo.ListByIDs(ctx, seatIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list seats by IDs")
		return nil, fmt.Errorf("failed to list seats: %w", err)
	}
	return seats, nil
}

func (s *SeatServiceImpl) CreateSeatsFromTemplate(ctx context.Context, req *model.BulkCreateSeatsRequest) ([]model.Seat, error) {
	// Verify bus exists
	_, err := s.busRepo.GetBusByID(ctx, req.BusID)
	if err != nil {
		log.Error().Err(err).Str("bus_id", req.BusID.String()).Msg("Bus not found")
		return nil, fmt.Errorf("bus not found: %w", err)
	}

	var seats []model.Seat
	for _, seatReq := range req.Seats {
		seat := model.Seat{
			BusID:       req.BusID,
			SeatNumber:  seatReq.SeatNumber,
			Row:         seatReq.Row,
			Column:      seatReq.Column,
			SeatType:    seatReq.SeatType,
			Floor:       seatReq.Floor,
			IsAvailable: true,
		}

		if seatReq.PriceMultiplier != nil {
			seat.PriceMultiplier = *seatReq.PriceMultiplier
		} else {
			seat.PriceMultiplier = seatReq.SeatType.GetPriceMultiplier()
		}

		seats = append(seats, seat)
	}

	if err := s.seatRepo.CreateBulk(ctx, seats); err != nil {
		log.Error().Err(err).Msg("Failed to create seats in bulk")
		return nil, fmt.Errorf("failed to create seats: %w", err)
	}

	log.Info().Str("bus_id", req.BusID.String()).Int("count", len(seats)).Msg("Seats created successfully")
	return seats, nil
}

func (s *SeatServiceImpl) GetSeatMap(ctx context.Context, busID uuid.UUID) (*model.SeatMapResponse, error) {
	seats, err := s.seatRepo.ListByBusID(ctx, busID)
	if err != nil {
		log.Error().Err(err).Str("bus_id", busID.String()).Msg("Failed to get seat map")
		return nil, fmt.Errorf("failed to get seat map: %w", err)
	}

	// Convert to seat details
	var seatDetails []model.SeatDetail
	maxRows, maxCols, maxFloors := 0, 0, 1

	for _, seat := range seats {
		seatDetails = append(seatDetails, model.SeatDetail{
			ID:              seat.ID,
			SeatNumber:      seat.SeatNumber,
			Row:             seat.Row,
			Column:          seat.Column,
			SeatType:        seat.SeatType.String(),
			PriceMultiplier: seat.PriceMultiplier,
			IsAvailable:     seat.IsAvailable,
			Floor:           seat.Floor,
		})

		if seat.Row > maxRows {
			maxRows = seat.Row
		}
		if seat.Column > maxCols {
			maxCols = seat.Column
		}
		if seat.Floor > maxFloors {
			maxFloors = seat.Floor
		}
	}

	response := &model.SeatMapResponse{
		BusID:      busID,
		TotalSeats: len(seats),
		Seats:      seatDetails,
		Layout: model.SeatLayoutInfo{
			MaxRows:    maxRows,
			MaxColumns: maxCols,
			Floors:     maxFloors,
		},
	}

	return response, nil
}

func (s *SeatServiceImpl) CreateSeat(ctx context.Context, req *model.CreateSeatRequest) (*model.Seat, error) {
	_, err := s.busRepo.GetBusByID(ctx, req.BusID)
	if err != nil {
		log.Error().Err(err).Str("bus_id", req.BusID.String()).Msg("Bus not found")
		return nil, ginext.NewBadRequestError("bus not found")
	}

	seat := &model.Seat{
		BusID:       req.BusID,
		SeatNumber:  req.SeatNumber,
		Row:         req.Row,
		Column:      req.Column,
		SeatType:    req.SeatType,
		Floor:       req.Floor,
		IsAvailable: true,
	}

	if req.PriceMultiplier != nil {
		seat.PriceMultiplier = *req.PriceMultiplier
	} else {
		seat.PriceMultiplier = req.SeatType.GetPriceMultiplier()
	}

	if err := s.seatRepo.Create(ctx, seat); err != nil {
		log.Error().Err(err).Msg("Failed to create seat")
		return nil, ginext.NewInternalServerError("failed to create seat")
	}
	return seat, nil
}

func (s *SeatServiceImpl) UpdateSeat(ctx context.Context, id uuid.UUID, req *model.UpdateSeatRequest) (*model.Seat, error) {
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

func (s *SeatServiceImpl) DeleteSeat(ctx context.Context, id uuid.UUID) error {
	if err := s.seatRepo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Str("seat_id", id.String()).Msg("Failed to delete seat")
		return ginext.NewInternalServerError("failed to delete seat")
	}
	return nil
}
