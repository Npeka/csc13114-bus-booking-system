package service

import (
	"context"
	"errors"
	"fmt"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BusService interface {
	CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error)
	GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error)
	DeleteBus(ctx context.Context, id uuid.UUID) error
	ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error)
	GetBusSeats(ctx context.Context, busID uuid.UUID) ([]model.Seat, error)
}

type BusServiceImpl struct {
	repositories *repository.Repositories
}

func NewBusService(repositories *repository.Repositories) BusService {
	return &BusServiceImpl{
		repositories: repositories,
	}
}

func (s *BusServiceImpl) CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error) {
	log.Info().Msg("Creating new bus")

	// Validate request
	if req.PlateNumber == "" || req.Model == "" {
		return nil, errors.New("plate number and model are required")
	}

	if req.SeatCapacity <= 0 || req.SeatCapacity > 100 {
		return nil, errors.New("seat capacity must be between 1 and 100")
	}

	// Check if plate number already exists
	existing, err := s.repositories.Bus.GetBusByPlateNumber(ctx, req.PlateNumber)
	if err == nil && existing != nil {
		return nil, errors.New("plate number already exists")
	}

	bus := &model.Bus{
		OperatorID:   req.OperatorID,
		PlateNumber:  req.PlateNumber,
		Model:        req.Model,
		SeatCapacity: req.SeatCapacity,
		Amenities:    req.Amenities,
		IsActive:     true,
	}

	if err := s.repositories.Bus.CreateBus(ctx, bus); err != nil {
		log.Error().Err(err).Msg("Failed to create bus")
		return nil, fmt.Errorf("failed to create bus: %w", err)
	}

	// Create seats for the bus
	if err := s.createSeatsForBus(ctx, bus); err != nil {
		log.Error().Err(err).Str("bus_id", bus.ID.String()).Msg("Failed to create seats for bus")
		// Continue anyway, seats can be created later
	}

	log.Info().Str("bus_id", bus.ID.String()).Msg("Bus created successfully")
	return bus, nil
}

func (s *BusServiceImpl) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	bus, err := s.repositories.Bus.GetBusByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to get bus")
		return nil, fmt.Errorf("failed to get bus: %w", err)
	}

	return bus, nil
}

func (s *BusServiceImpl) UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error) {
	log.Info().Str("bus_id", id.String()).Msg("Updating bus")

	// Get existing bus
	bus, err := s.repositories.Bus.GetBusByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("bus not found: %w", err)
	}

	// Update fields
	if req.PlateNumber != nil {
		// Check if new plate number already exists (and it's not the same bus)
		existing, err := s.repositories.Bus.GetBusByPlateNumber(ctx, *req.PlateNumber)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("plate number already exists")
		}
		bus.PlateNumber = *req.PlateNumber
	}
	if req.Model != nil {
		bus.Model = *req.Model
	}
	if req.SeatCapacity != nil {
		if *req.SeatCapacity <= 0 || *req.SeatCapacity > 100 {
			return nil, errors.New("seat capacity must be between 1 and 100")
		}
		bus.SeatCapacity = *req.SeatCapacity
	}
	if req.Amenities != nil {
		bus.Amenities = *req.Amenities
	}
	if req.IsActive != nil {
		bus.IsActive = *req.IsActive
	}

	if err := s.repositories.Bus.UpdateBus(ctx, bus); err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to update bus")
		return nil, fmt.Errorf("failed to update bus: %w", err)
	}

	log.Info().Str("bus_id", id.String()).Msg("Bus updated successfully")
	return bus, nil
}

func (s *BusServiceImpl) DeleteBus(ctx context.Context, id uuid.UUID) error {
	log.Info().Str("bus_id", id.String()).Msg("Deleting bus")

	if err := s.repositories.Bus.DeleteBus(ctx, id); err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to delete bus")
		return fmt.Errorf("failed to delete bus: %w", err)
	}

	log.Info().Str("bus_id", id.String()).Msg("Bus deleted successfully")
	return nil
}

func (s *BusServiceImpl) ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	buses, total, err := s.repositories.Bus.ListBuses(ctx, operatorID, page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list buses")
		return nil, 0, fmt.Errorf("failed to list buses: %w", err)
	}

	return buses, total, nil
}

func (s *BusServiceImpl) GetBusSeats(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	seats, err := s.repositories.Seat.GetSeatsByBusID(ctx, busID)
	if err != nil {
		log.Error().Err(err).Str("bus_id", busID.String()).Msg("Failed to get bus seats")
		return nil, fmt.Errorf("failed to get bus seats: %w", err)
	}

	return seats, nil
}

// createSeatsForBus creates standard seats for a new bus
func (s *BusServiceImpl) createSeatsForBus(ctx context.Context, bus *model.Bus) error {
	seats := make([]model.Seat, 0, bus.SeatCapacity)

	// Create seats with standard naming (A1, A2, B1, B2, etc.)
	rowNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	seatsPerRow := 4
	if bus.SeatCapacity > 40 {
		seatsPerRow = 5
	}

	seatCount := 0
	for rowIdx := 0; rowIdx < len(rowNames) && seatCount < bus.SeatCapacity; rowIdx++ {
		for seatNum := 1; seatNum <= seatsPerRow && seatCount < bus.SeatCapacity; seatNum++ {
			seatCode := fmt.Sprintf("%s%d", rowNames[rowIdx], seatNum)
			seatType := "standard"

			// First and last rows are premium
			if rowIdx == 0 || rowIdx == len(rowNames)-1 {
				seatType = "premium"
			}

			seat := model.Seat{
				BusID:    bus.ID,
				SeatCode: seatCode,
				SeatType: seatType,
				IsActive: true,
			}
			seats = append(seats, seat)
			seatCount++
		}
	}

	return s.repositories.Seat.CreateSeats(ctx, seats)
}
