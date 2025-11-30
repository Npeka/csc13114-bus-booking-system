package service

import (
	"context"
	"fmt"
	"time"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TripService interface {
	SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error)
	GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error)
	ListTrips(ctx context.Context, page, pageSize int) ([]model.Trip, int64, error)
	GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error)
	GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, departureDate time.Time) ([]model.Trip, error)

	CreateTrip(ctx context.Context, req *model.CreateTripRequest) (*model.Trip, error)
	UpdateTrip(ctx context.Context, id uuid.UUID, req *model.UpdateTripRequest) (*model.Trip, error)
	DeleteTrip(ctx context.Context, id uuid.UUID) error
}

type TripServiceImpl struct {
	tripRepo      repository.TripRepository
	routeRepo     repository.RouteRepository
	routeStopRepo repository.RouteStopRepository
	busRepo       repository.BusRepository
	seatRepo      repository.SeatRepository
}

func NewTripService(
	tripRepo repository.TripRepository,
	routeRepo repository.RouteRepository,
	routeStopRepo repository.RouteStopRepository,
	busRepo repository.BusRepository,
	seatRepo repository.SeatRepository,
) TripService {
	return &TripServiceImpl{
		tripRepo:      tripRepo,
		routeRepo:     routeRepo,
		routeStopRepo: routeStopRepo,
		busRepo:       busRepo,
		seatRepo:      seatRepo,
	}
}

func (s *TripServiceImpl) SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error) {
	date, err := time.Parse("02/01/2006", req.Date)
	if err != nil {
		return nil, 0, ginext.NewBadRequestError("invalid date format, use DD/MM/YYYY")
	}

	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, 0, ginext.NewBadRequestError("search date cannot be in the past")
	}

	trips, total, err := s.tripRepo.SearchTrips(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search trips")
		return nil, 0, ginext.NewInternalServerError("failed to search trips")
	}

	return trips, total, nil
}

func (s *TripServiceImpl) GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get trip")
	}
	return trip, nil
}

func (s *TripServiceImpl) ListTrips(ctx context.Context, page, pageSize int) ([]model.Trip, int64, error) {
	trips, total, err := s.tripRepo.ListTrips(ctx, page, pageSize)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("failed to list trips")
	}
	return trips, total, nil
}

func (s *TripServiceImpl) GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, ginext.NewInternalServerError("trip not found")
	}

	seats, err := s.seatRepo.ListByBusID(ctx, trip.BusID)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get seats")
	}

	// TODO: Check seat status from booking service
	var seatAvailabilities []model.SeatAvailability
	availableCount := 0

	for _, seat := range seats {
		seatAvail := model.SeatAvailability{
			SeatID:      seat.ID,
			SeatNumber:  seat.SeatNumber,
			SeatType:    seat.SeatType,
			Price:       trip.BasePrice * seat.PriceMultiplier,
			IsAvailable: seat.IsAvailable, // TODO: Check from booking service
			Row:         seat.Row,
			Column:      seat.Column,
			Floor:       seat.Floor,
		}

		if seatAvail.IsAvailable {
			availableCount++
		}

		seatAvailabilities = append(seatAvailabilities, seatAvail)
	}

	return &model.SeatAvailabilityResponse{
		TripID:         tripID,
		AvailableSeats: availableCount,
		TotalSeats:     len(seats),
		SeatMap:        seatAvailabilities,
	}, nil
}

// GetTripsByRouteAndDate gets trips by route and departure date
func (s *TripServiceImpl) GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, departureDate time.Time) ([]model.Trip, error) {
	// Validate inputs
	if routeID == uuid.Nil {
		return nil, ginext.NewBadRequestError("route ID is required")
	}

	// Check if route exists
	_, err := s.routeRepo.GetRouteByID(ctx, routeID)
	if err != nil {
		return nil, fmt.Errorf("invalid route: %w", err)
	}

	// Get trips by route and date
	trips, err := s.tripRepo.GetTripsByRouteAndDate(ctx, routeID, departureDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips: %w", err)
	}

	return trips, nil
}

func (s *TripServiceImpl) CreateTrip(ctx context.Context, req *model.CreateTripRequest) (*model.Trip, error) {
	if req.ArrivalTime.Before(req.DepartureTime) {
		return nil, ginext.NewBadRequestError("arrival time must be after departure time")
	}

	if req.DepartureTime.Before(time.Now()) {
		return nil, ginext.NewBadRequestError("departure time cannot be in the past")
	}

	// Check if route exists
	_, err := s.routeRepo.GetRouteByID(ctx, req.RouteID)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid route")
	}

	// Check if bus exists and is available
	bus, err := s.busRepo.GetBusByID(ctx, req.BusID)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid bus")
	}

	if !bus.IsActive {
		return nil, ginext.NewBadRequestError("bus is not active")
	}

	// Check for bus conflicts (same bus cannot have overlapping trips)
	conflictTrips, err := s.tripRepo.GetTripsByBusAndDateRange(ctx, req.BusID,
		req.DepartureTime.Add(-4*time.Hour), req.ArrivalTime.Add(4*time.Hour))
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to check bus availability")
	}

	for _, existingTrip := range conflictTrips {
		if req.ArrivalTime.After(existingTrip.DepartureTime) && req.DepartureTime.Before(existingTrip.ArrivalTime) {
			return nil, ginext.NewBadRequestError("bus is already assigned to another trip during the specified time")
		}
	}

	trip := &model.Trip{
		RouteID:       req.RouteID,
		BusID:         req.BusID,
		DepartureTime: req.DepartureTime,
		ArrivalTime:   req.ArrivalTime,
		BasePrice:     req.BasePrice,
		Status:        "scheduled",
		IsActive:      true,
	}

	if err := s.tripRepo.CreateTrip(ctx, trip); err != nil {
		log.Error().Err(err).Msg("Failed to create trip")
		return nil, ginext.NewInternalServerError("failed to create trip")
	}

	// Load relationships
	return s.GetTripByID(ctx, trip.ID)
}

func (s *TripServiceImpl) UpdateTrip(ctx context.Context, id uuid.UUID, req *model.UpdateTripRequest) (*model.Trip, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get trip")
	}

	// Update fields if provided
	if req.DepartureTime != nil {
		if req.DepartureTime.Before(time.Now()) {
			return nil, ginext.NewBadRequestError("departure time cannot be in the past")
		}
		trip.DepartureTime = *req.DepartureTime
	}

	if req.ArrivalTime != nil {
		if req.ArrivalTime.Before(trip.DepartureTime) {
			return nil, ginext.NewBadRequestError("arrival time must be after departure time")
		}
		trip.ArrivalTime = *req.ArrivalTime
	}

	if req.BasePrice != nil {
		if *req.BasePrice < 0 {
			return nil, ginext.NewBadRequestError("base price must be non-negative")
		}
		trip.BasePrice = *req.BasePrice
	}

	if req.Status != nil {
		trip.Status = *req.Status
	}

	if req.IsActive != nil {
		trip.IsActive = *req.IsActive
	}

	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return nil, ginext.NewInternalServerError("failed to update trip")
	}

	return s.GetTripByID(ctx, id)
}

func (s *TripServiceImpl) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		return fmt.Errorf("trip not found: %w", err)
	}

	if trip.Status != "scheduled" {
		return ginext.NewBadRequestError("only scheduled trips can be deleted")
	}

	if trip.DepartureTime.Before(time.Now().Add(24 * time.Hour)) {
		return ginext.NewBadRequestError("cannot delete trip within 24 hours of departure")
	}

	if err := s.tripRepo.DeleteTrip(ctx, id); err != nil {
		return ginext.NewInternalServerError("failed to delete trip")
	}

	return nil
}
