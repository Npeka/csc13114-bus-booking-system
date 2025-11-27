package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TripService interface {
	// Trip operations
	SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error)
	GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error)
	CreateTrip(ctx context.Context, req *model.CreateTripRequest) (*model.Trip, error)
	UpdateTrip(ctx context.Context, id uuid.UUID, req *model.UpdateTripRequest) (*model.Trip, error)
	DeleteTrip(ctx context.Context, id uuid.UUID) error
	GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error)
	GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, departureDate time.Time) ([]model.Trip, error)

	// Route operations
	ListRoutes(ctx context.Context, operatorID *uuid.UUID, page, limit int) (*model.RouteListResponse, error)
	GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error)
	CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (*model.Route, error)
	UpdateRoute(ctx context.Context, id uuid.UUID, req *model.UpdateRouteRequest) (*model.Route, error)
	DeleteRoute(ctx context.Context, id uuid.UUID) error

	// Bus operations
	ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error)
	GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error)
	UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error)
	DeleteBus(ctx context.Context, id uuid.UUID) error

	// Operator operations
	ListOperators(ctx context.Context, page, limit int) (*model.OperatorListResponse, error)
	GetOperatorByID(ctx context.Context, id uuid.UUID) (*model.Operator, error)
	CreateOperator(ctx context.Context, req *model.CreateOperatorRequest) (*model.Operator, error)
	UpdateOperator(ctx context.Context, id uuid.UUID, req *model.UpdateOperatorRequest) (*model.Operator, error)
	DeleteOperator(ctx context.Context, id uuid.UUID) error
}

type TripServiceImpl struct {
	tripRepo      repository.TripRepository
	routeRepo     repository.RouteRepository
	routeStopRepo repository.RouteStopRepository
	busRepo       repository.BusRepository
	operatorRepo  repository.OperatorRepository
	seatRepo      repository.SeatRepository
}

func NewTripService(
	tripRepo repository.TripRepository,
	routeRepo repository.RouteRepository,
	routeStopRepo repository.RouteStopRepository,
	busRepo repository.BusRepository,
	operatorRepo repository.OperatorRepository,
	seatRepo repository.SeatRepository,
) TripService {
	return &TripServiceImpl{
		tripRepo:      tripRepo,
		routeRepo:     routeRepo,
		routeStopRepo: routeStopRepo,
		busRepo:       busRepo,
		operatorRepo:  operatorRepo,
		seatRepo:      seatRepo,
	}
}

// Trip operations
func (s *TripServiceImpl) SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error) {
	// Validate search request
	if req.Origin == "" || req.Destination == "" {
		return nil, 0, errors.New("origin and destination are required")
	}

	if req.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, 0, errors.New("departure date cannot be in the past")
	}

	// Set default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	trips, total, err := s.tripRepo.SearchTrips(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search trips")
		return nil, 0, fmt.Errorf("failed to search trips: %w", err)
	}

	return trips, total, nil
}

func (s *TripServiceImpl) GetTripByID(ctx context.Context, id uuid.UUID) (*model.Trip, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("trip_id", id.String()).Msg("Failed to get trip by ID")
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}
	return trip, nil
}

func (s *TripServiceImpl) CreateTrip(ctx context.Context, req *model.CreateTripRequest) (*model.Trip, error) {
	// Validate trip times
	if req.ArrivalTime.Before(req.DepartureTime) {
		return nil, errors.New("arrival time must be after departure time")
	}

	if req.DepartureTime.Before(time.Now()) {
		return nil, errors.New("departure time cannot be in the past")
	}

	// Check if route exists
	_, err := s.routeRepo.GetRouteByID(ctx, req.RouteID)
	if err != nil {
		return nil, fmt.Errorf("invalid route: %w", err)
	}

	// Check if bus exists and is available
	bus, err := s.busRepo.GetBusByID(ctx, req.BusID)
	if err != nil {
		return nil, fmt.Errorf("invalid bus: %w", err)
	}

	if !bus.IsActive {
		return nil, errors.New("bus is not active")
	}

	// Check for bus conflicts (same bus cannot have overlapping trips)
	conflictTrips, err := s.tripRepo.GetTripsByBusAndDateRange(ctx, req.BusID,
		req.DepartureTime.Add(-4*time.Hour), req.ArrivalTime.Add(4*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to check bus availability: %w", err)
	}

	for _, existingTrip := range conflictTrips {
		if req.ArrivalTime.After(existingTrip.DepartureTime) && req.DepartureTime.Before(existingTrip.ArrivalTime) {
			return nil, errors.New("bus is not available at the specified time")
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
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	// Load relationships
	return s.GetTripByID(ctx, trip.ID)
}

func (s *TripServiceImpl) UpdateTrip(ctx context.Context, id uuid.UUID, req *model.UpdateTripRequest) (*model.Trip, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("trip not found: %w", err)
	}

	// Update fields if provided
	if req.DepartureTime != nil {
		if req.DepartureTime.Before(time.Now()) {
			return nil, errors.New("departure time cannot be in the past")
		}
		trip.DepartureTime = *req.DepartureTime
	}

	if req.ArrivalTime != nil {
		if req.ArrivalTime.Before(trip.DepartureTime) {
			return nil, errors.New("arrival time must be after departure time")
		}
		trip.ArrivalTime = *req.ArrivalTime
	}

	if req.BasePrice != nil {
		if *req.BasePrice < 0 {
			return nil, errors.New("base price must be non-negative")
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
		log.Error().Err(err).Str("trip_id", id.String()).Msg("Failed to update trip")
		return nil, fmt.Errorf("failed to update trip: %w", err)
	}

	return s.GetTripByID(ctx, id)
}

func (s *TripServiceImpl) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		return fmt.Errorf("trip not found: %w", err)
	}

	// Check if trip can be deleted (not started and no bookings)
	if trip.Status != "scheduled" {
		return errors.New("cannot delete trip that is not scheduled")
	}

	if trip.DepartureTime.Before(time.Now().Add(24 * time.Hour)) {
		return errors.New("cannot delete trip within 24 hours of departure")
	}

	if err := s.tripRepo.DeleteTrip(ctx, id); err != nil {
		log.Error().Err(err).Str("trip_id", id.String()).Msg("Failed to delete trip")
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	return nil
}

func (s *TripServiceImpl) GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("trip not found: %w", err)
	}

	seats, err := s.seatRepo.ListByBusID(ctx, trip.BusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
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

// Route operations
func (s *TripServiceImpl) ListRoutes(ctx context.Context, operatorID *uuid.UUID, page, limit int) (*model.RouteListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	routes, total, err := s.routeRepo.ListRoutes(ctx, operatorID, page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list routes")
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &model.RouteListResponse{
		Routes:     routes,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *TripServiceImpl) GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	route, err := s.routeRepo.GetRouteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}
	return route, nil
}

func (s *TripServiceImpl) CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (*model.Route, error) {
	// Validate operator exists
	_, err := s.operatorRepo.GetOperatorByID(ctx, req.OperatorID)
	if err != nil {
		return nil, fmt.Errorf("invalid operator: %w", err)
	}

	if req.Origin == req.Destination {
		return nil, errors.New("origin and destination must be different")
	}

	route := &model.Route{
		OperatorID:       req.OperatorID,
		Origin:           req.Origin,
		Destination:      req.Destination,
		DistanceKm:       req.DistanceKm,
		EstimatedMinutes: req.EstimatedMinutes,
		IsActive:         true,
	}

	if err := s.routeRepo.CreateRoute(ctx, route); err != nil {
		log.Error().Err(err).Msg("Failed to create route")
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	return s.GetRouteByID(ctx, route.ID)
}

func (s *TripServiceImpl) UpdateRoute(ctx context.Context, id uuid.UUID, req *model.UpdateRouteRequest) (*model.Route, error) {
	route, err := s.routeRepo.GetRouteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	// Update fields if provided
	if req.Origin != nil {
		route.Origin = *req.Origin
	}
	if req.Destination != nil {
		route.Destination = *req.Destination
	}
	if req.DistanceKm != nil {
		route.DistanceKm = *req.DistanceKm
	}
	if req.EstimatedMinutes != nil {
		route.EstimatedMinutes = *req.EstimatedMinutes
	}
	if req.IsActive != nil {
		route.IsActive = *req.IsActive
	}

	// Validate origin != destination
	if route.Origin == route.Destination {
		return nil, errors.New("origin and destination must be different")
	}

	if err := s.routeRepo.UpdateRoute(ctx, route); err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to update route")
		return nil, fmt.Errorf("failed to update route: %w", err)
	}

	return s.GetRouteByID(ctx, id)
}

func (s *TripServiceImpl) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	// Check if route has active trips
	_, err := s.routeRepo.GetRouteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("route not found: %w", err)
	}

	// TODO: Check if route has active trips in the future

	if err := s.routeRepo.DeleteRoute(ctx, id); err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to delete route")
		return fmt.Errorf("failed to delete route: %w", err)
	}

	return nil
}

// Bus operations
func (s *TripServiceImpl) ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	buses, total, err := s.busRepo.ListBuses(ctx, operatorID, page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list buses")
		return nil, 0, fmt.Errorf("failed to list buses: %w", err)
	}

	return buses, total, nil
}

func (s *TripServiceImpl) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	bus, err := s.busRepo.GetBusByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("bus not found: %w", err)
	}
	return bus, nil
}

func (s *TripServiceImpl) CreateBus(ctx context.Context, req *model.CreateBusRequest) (*model.Bus, error) {
	// Validate operator exists
	_, err := s.operatorRepo.GetOperatorByID(ctx, req.OperatorID)
	if err != nil {
		return nil, fmt.Errorf("invalid operator: %w", err)
	}

	// Check plate number uniqueness
	existingBus, err := s.busRepo.GetBusByPlateNumber(ctx, req.PlateNumber)
	if err == nil && existingBus != nil {
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

	if err := s.busRepo.CreateBus(ctx, bus); err != nil {
		log.Error().Err(err).Msg("Failed to create bus")
		return nil, fmt.Errorf("failed to create bus: %w", err)
	}

	// Generate seats for the bus
	// Note: Seat generation is now handled by BusService
	// if err := s.generateSeatsForBus(ctx, bus.ID, req.SeatCapacity); err != nil {
	// 	log.Error().Err(err).Str("bus_id", bus.ID.String()).Msg("Failed to generate seats for bus")
	// }

	return s.GetBusByID(ctx, bus.ID)
}

func (s *TripServiceImpl) UpdateBus(ctx context.Context, id uuid.UUID, req *model.UpdateBusRequest) (*model.Bus, error) {
	bus, err := s.busRepo.GetBusByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("bus not found: %w", err)
	}

	// Update fields if provided
	if req.PlateNumber != nil {
		// Check uniqueness if changing plate number
		if *req.PlateNumber != bus.PlateNumber {
			existingBus, err := s.busRepo.GetBusByPlateNumber(ctx, *req.PlateNumber)
			if err == nil && existingBus != nil {
				return nil, errors.New("plate number already exists")
			}
		}
		bus.PlateNumber = *req.PlateNumber
	}

	if req.Model != nil {
		bus.Model = *req.Model
	}

	if req.SeatCapacity != nil {
		bus.SeatCapacity = *req.SeatCapacity
		// TODO: Handle seat capacity changes (may need to update seats)
	}

	if req.Amenities != nil {
		bus.Amenities = *req.Amenities
	}

	if req.IsActive != nil {
		bus.IsActive = *req.IsActive
	}

	if err := s.busRepo.UpdateBus(ctx, bus); err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to update bus")
		return nil, fmt.Errorf("failed to update bus: %w", err)
	}

	return s.GetBusByID(ctx, id)
}

func (s *TripServiceImpl) DeleteBus(ctx context.Context, id uuid.UUID) error {
	// TODO: Check if bus has active trips

	if err := s.busRepo.DeleteBus(ctx, id); err != nil {
		log.Error().Err(err).Str("bus_id", id.String()).Msg("Failed to delete bus")
		return fmt.Errorf("failed to delete bus: %w", err)
	}

	return nil
}

// Operator operations
func (s *TripServiceImpl) ListOperators(ctx context.Context, page, limit int) (*model.OperatorListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	operators, total, err := s.operatorRepo.ListOperators(ctx, page, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list operators")
		return nil, fmt.Errorf("failed to list operators: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &model.OperatorListResponse{
		Operators:  operators,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *TripServiceImpl) GetOperatorByID(ctx context.Context, id uuid.UUID) (*model.Operator, error) {
	operator, err := s.operatorRepo.GetOperatorByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("operator not found: %w", err)
	}
	return operator, nil
}

func (s *TripServiceImpl) CreateOperator(ctx context.Context, req *model.CreateOperatorRequest) (*model.Operator, error) {
	// Check email uniqueness
	existingOperator, err := s.operatorRepo.GetOperatorByEmail(ctx, req.ContactEmail)
	if err == nil && existingOperator != nil {
		return nil, errors.New("email already exists")
	}

	operator := &model.Operator{
		Name:         req.Name,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Status:       "pending",
	}

	if err := s.operatorRepo.CreateOperator(ctx, operator); err != nil {
		log.Error().Err(err).Msg("Failed to create operator")
		return nil, fmt.Errorf("failed to create operator: %w", err)
	}

	return s.GetOperatorByID(ctx, operator.ID)
}

func (s *TripServiceImpl) UpdateOperator(ctx context.Context, id uuid.UUID, req *model.UpdateOperatorRequest) (*model.Operator, error) {
	operator, err := s.operatorRepo.GetOperatorByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("operator not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		operator.Name = *req.Name
	}

	if req.ContactEmail != nil {
		// Check uniqueness if changing email
		if *req.ContactEmail != operator.ContactEmail {
			existingOperator, err := s.operatorRepo.GetOperatorByEmail(ctx, *req.ContactEmail)
			if err == nil && existingOperator != nil {
				return nil, errors.New("email already exists")
			}
		}
		operator.ContactEmail = *req.ContactEmail
	}

	if req.ContactPhone != nil {
		operator.ContactPhone = *req.ContactPhone
	}

	if err := s.operatorRepo.UpdateOperator(ctx, operator); err != nil {
		log.Error().Err(err).Str("operator_id", id.String()).Msg("Failed to update operator")
		return nil, fmt.Errorf("failed to update operator: %w", err)
	}

	return s.GetOperatorByID(ctx, id)
}

func (s *TripServiceImpl) DeleteOperator(ctx context.Context, id uuid.UUID) error {
	// TODO: Check if operator has active routes/buses/trips

	if err := s.operatorRepo.DeleteOperator(ctx, id); err != nil {
		log.Error().Err(err).Str("operator_id", id.String()).Msg("Failed to delete operator")
		return fmt.Errorf("failed to delete operator: %w", err)
	}

	return nil
}

// Helper functions

// GetTripsByRouteAndDate gets trips by route and departure date
func (s *TripServiceImpl) GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, departureDate time.Time) ([]model.Trip, error) {
	// Validate inputs
	if routeID == uuid.Nil {
		return nil, errors.New("route ID is required")
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
