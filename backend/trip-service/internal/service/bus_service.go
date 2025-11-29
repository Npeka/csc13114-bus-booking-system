package service

import (
	"context"
	"fmt"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BusService interface {
	GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	ListBuses(ctx context.Context, req model.ListBusesRequest) ([]model.Bus, int64, error)
	GetBusSeats(ctx context.Context, busID uuid.UUID) ([]model.Seat, error)

	CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error)
	UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error)
	DeleteBus(ctx context.Context, id uuid.UUID) error
}

type BusServiceImpl struct {
	busRepo  repository.BusRepository
	seatRepo repository.SeatRepository
}

func NewBusService(
	busRepo repository.BusRepository,
	seatRepo repository.SeatRepository,
) BusService {
	return &BusServiceImpl{
		busRepo:  busRepo,
		seatRepo: seatRepo,
	}
}

func (s *BusServiceImpl) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	bus, err := s.busRepo.GetBusWithSeatsByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to get bus")
		return nil, ginext.NewInternalServerError("failed to get bus")
	}
	return bus, nil
}

func (s *BusServiceImpl) ListBuses(ctx context.Context, req model.ListBusesRequest) ([]model.Bus, int64, error) {
	buses, total, err := s.busRepo.ListBuses(ctx, req.Page, req.PageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list buses")
		return nil, 0, ginext.NewInternalServerError("failed to list buses")
	}

	return buses, total, nil
}

func (s *BusServiceImpl) GetBusSeats(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	seats, err := s.seatRepo.ListByBusID(ctx, busID)
	if err != nil {
		log.Error().Err(err).Str("bus_id", busID.String()).Msg("Failed to get bus seats")
		return nil, fmt.Errorf("failed to get bus seats: %w", err)
	}

	return seats, nil
}

func (s *BusServiceImpl) CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error) {
	existing, err := s.busRepo.GetBusByPlateNumber(ctx, req.PlateNumber)
	if err == nil && existing != nil {
		return nil, ginext.NewBadRequestError("plate number already exists")
	}

	bus := &model.Bus{
		PlateNumber:  req.PlateNumber,
		Model:        req.Model,
		SeatCapacity: req.SeatCapacity,
		Amenities:    req.Amenities,
		IsActive:     true,
		Seats:        s.generateSeatsForBus(req.SeatCapacity),
	}

	if err := s.busRepo.CreateBus(ctx, bus); err != nil {
		log.Error().Err(err).Msg("Failed to create bus")
		return nil, ginext.NewInternalServerError("failed to create bus")
	}

	return bus, nil
}

func (s *BusServiceImpl) generateSeatsForBus(seatCapacity int) []model.Seat {
	seats := make([]model.Seat, 0, seatCapacity)

	rowNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}
	seatsPerRow := 4
	if seatCapacity > 40 {
		seatsPerRow = 5
	}

	seatCount := 0
	for rowIdx := 0; rowIdx < len(rowNames) && seatCount < seatCapacity; rowIdx++ {
		for seatNum := 1; seatNum <= seatsPerRow && seatCount < seatCapacity; seatNum++ {
			seatNumber := fmt.Sprintf("%s%d", rowNames[rowIdx], seatNum)
			seatType := constants.SeatTypeStandard

			// First and last rows are premium (VIP)
			if rowIdx == 0 || rowIdx == len(rowNames)-1 {
				seatType = constants.SeatTypeVIP
			}

			seat := model.Seat{
				SeatNumber:  seatNumber,
				SeatType:    seatType,
				Row:         rowIdx + 1,
				Column:      seatNum,
				IsAvailable: true,
			}
			seats = append(seats, seat)
			seatCount++
		}
	}

	return seats
}

func (s *BusServiceImpl) UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error) {
	bus, err := s.busRepo.GetBusByID(ctx, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get bus")
	}

	if req.PlateNumber != nil {
		existing, err := s.busRepo.GetBusByPlateNumber(ctx, *req.PlateNumber)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ginext.NewBadRequestError("plate number already exists")
		}
		bus.PlateNumber = *req.PlateNumber
	}

	if req.Model != nil {
		bus.Model = *req.Model
	}

	if req.SeatCapacity != nil {
		if *req.SeatCapacity <= 0 || *req.SeatCapacity > 100 {
			return nil, ginext.NewBadRequestError("seat capacity must be between 1 and 100")
		}
		bus.SeatCapacity = *req.SeatCapacity
	}

	if req.Amenities != nil {
		bus.Amenities = *req.Amenities
	}

	if req.IsActive != nil {
		bus.IsActive = *req.IsActive
	}

	if err := s.busRepo.UpdateBus(ctx, bus); err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to update bus")
		return nil, ginext.NewInternalServerError("failed to update bus")
	}

	return bus, nil
}

func (s *BusServiceImpl) DeleteBus(ctx context.Context, id uuid.UUID) error {
	if err := s.busRepo.DeleteBus(ctx, id); err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to delete bus")
		return ginext.NewInternalServerError("failed to delete bus")
	}
	return nil
}
