package service

import (
	"context"
	"fmt"
	"time"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/client"
	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/model/booking"
	"bus-booking/trip-service/internal/model/payment"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TripService interface {
	SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error)
	GetTripByID(ctx context.Context, req *model.GetTripByIDRequest, id uuid.UUID) (*model.Trip, error)
	ListTrips(ctx context.Context, req *model.ListTripsRequest) ([]model.Trip, int64, error)
	GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error)
	GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, departureDate time.Time) ([]model.Trip, error)
	GetCompletedTripsForReschedule(ctx context.Context) ([]model.Trip, error)

	CreateTrip(ctx context.Context, req *model.CreateTripRequest) (*model.Trip, error)
	UpdateTrip(ctx context.Context, id uuid.UUID, req *model.UpdateTripRequest) (*model.Trip, error)
	DeleteTrip(ctx context.Context, id uuid.UUID) error
	RescheduleTrip(ctx context.Context, id uuid.UUID, newDeparture, newArrival time.Time) error
	CancelTrip(ctx context.Context, id uuid.UUID) error
	ProcessTripStatusUpdates(ctx context.Context) error
}

type TripServiceImpl struct {
	tripRepo      repository.TripRepository
	routeRepo     repository.RouteRepository
	routeStopRepo repository.RouteStopRepository
	busRepo       repository.BusRepository
	seatRepo      repository.SeatRepository
	bookingClient client.BookingClient
	paymentClient client.PaymentClient
}

func NewTripService(
	tripRepo repository.TripRepository,
	routeRepo repository.RouteRepository,
	routeStopRepo repository.RouteStopRepository,
	busRepo repository.BusRepository,
	seatRepo repository.SeatRepository,
	bookingClient client.BookingClient,
	paymentClient client.PaymentClient,
) TripService {
	return &TripServiceImpl{
		tripRepo:      tripRepo,
		routeRepo:     routeRepo,
		routeStopRepo: routeStopRepo,
		busRepo:       busRepo,
		seatRepo:      seatRepo,
		bookingClient: bookingClient,
		paymentClient: paymentClient,
	}
}

func (s *TripServiceImpl) SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error) {
	trips, total, err := s.tripRepo.SearchTrips(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search trips")
		return nil, 0, ginext.NewInternalServerError("failed to search trips")
	}

	return trips, total, nil
}

func (s *TripServiceImpl) GetTripByID(ctx context.Context, req *model.GetTripByIDRequest, id uuid.UUID) (*model.Trip, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, req, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get trip")
	}

	if req.SeatBookingStatus && req.PreloadBus && req.PreloadSeat && trip.Bus != nil && len(trip.Bus.Seats) > 0 {
		seatIDs := make([]uuid.UUID, len(trip.Bus.Seats))
		for i, seat := range trip.Bus.Seats {
			seatIDs[i] = seat.ID
		}

		seatStatuses, err := s.bookingClient.GetSeatStatus(ctx, trip.ID, seatIDs)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check seat status from booking service")
			return nil, ginext.NewInternalServerError("failed to get seat status")
		}

		seatStatusMap := make(map[uuid.UUID]booking.SeatStatus)
		for _, status := range seatStatuses {
			seatStatusMap[status.SeatID] = status
		}

		for i, seat := range trip.Bus.Seats {
			if status, ok := seatStatusMap[seat.ID]; ok {
				trip.Bus.Seats[i].Status = &status
				continue
			}
			trip.Bus.Seats[i].Status = &booking.SeatStatus{
				SeatID:   seat.ID,
				IsBooked: false,
				IsLocked: false,
			}
		}
	}

	return trip, nil
}

func (s *TripServiceImpl) ListTrips(ctx context.Context, req *model.ListTripsRequest) ([]model.Trip, int64, error) {
	// If IDs provided, fetch specific trips (batch mode)
	if len(req.IDs) > 0 {
		trips, err := s.tripRepo.GetTripsByIDs(ctx, req.IDs)
		if err != nil {
			return nil, 0, ginext.NewInternalServerError("failed to list trips by IDs")
		}
		return trips, int64(len(trips)), nil
	}

	// Otherwise, use pagination
	trips, total, err := s.tripRepo.ListTrips(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("failed to list trips")
	}
	return trips, total, nil
}

func (s *TripServiceImpl) GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error) {
	req := &model.GetTripByIDRequest{
		PreloadBus:  true,
		PreloadSeat: true,
	}
	trip, err := s.tripRepo.GetTripByID(ctx, req, tripID)
	if err != nil {
		return nil, ginext.NewInternalServerError("trip not found")
	}

	seats, err := s.seatRepo.GetListByBusID(ctx, trip.BusID)
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
	return s.GetTripByID(ctx, &model.GetTripByIDRequest{}, trip.ID)
}

func (s *TripServiceImpl) UpdateTrip(ctx context.Context, id uuid.UUID, req *model.UpdateTripRequest) (*model.Trip, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, &model.GetTripByIDRequest{}, id)
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

	return s.GetTripByID(ctx, &model.GetTripByIDRequest{}, id)
}

func (s *TripServiceImpl) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	trip, err := s.tripRepo.GetTripByID(ctx, &model.GetTripByIDRequest{}, id)
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

// GetCompletedTripsForReschedule gets completed trips from last 24h that need rescheduling
func (s *TripServiceImpl) GetCompletedTripsForReschedule(ctx context.Context) ([]model.Trip, error) {
	return s.tripRepo.GetCompletedTripsForReschedule(ctx)
}

// RescheduleTrip updates trip with new departure/arrival times (maintains same bus, route, price)
func (s *TripServiceImpl) RescheduleTrip(ctx context.Context, id uuid.UUID, newDeparture, newArrival time.Time) error {
	trip, err := s.tripRepo.GetTripByID(ctx, &model.GetTripByIDRequest{}, id)
	if err != nil {
		return fmt.Errorf("trip not found: %w", err)
	}

	if newArrival.Before(newDeparture) {
		return ginext.NewBadRequestError("arrival time must be after departure time")
	}

	// Update times and reset status to scheduled
	trip.DepartureTime = newDeparture
	trip.ArrivalTime = newArrival
	trip.Status = "scheduled"

	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return fmt.Errorf("failed to reschedule trip: %w", err)
	}

	return nil
}

// ProcessTripStatusUpdates triggers the batch update of trip statuses
func (s *TripServiceImpl) ProcessTripStatusUpdates(ctx context.Context) error {
	return s.tripRepo.UpdateTripStatuses(ctx)
}

func (s *TripServiceImpl) CancelTrip(ctx context.Context, id uuid.UUID) error {
	// 1. Get Trip
	trip, err := s.tripRepo.GetTripByID(ctx, &model.GetTripByIDRequest{}, id)
	if err != nil {
		return ginext.NewInternalServerError("failed to get trip")
	}

	// 2. Validate Status Check
	// Only Scheduled or Delayed trips can be cancelled
	if trip.Status != constants.TripStatusScheduled && trip.Status != constants.TripStatusDelayed {
		return ginext.NewBadRequestError(fmt.Sprintf("Cannot cancel trip with status: %s. Only scheduled or delayed trips can be cancelled.", trip.Status))
	}

	// 3. Update Status to Cancelled
	trip.Status = constants.TripStatusCancelled
	
	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return ginext.NewInternalServerError("failed to update trip status")
	}

	// 4. Get Bookings
	bookings, err := s.bookingClient.GetTripBookings(ctx, id)
	if err != nil {
		// Log error but we already cancelled the trip in DB. 
		log.Error().Err(err).Str("trip_id", id.String()).Msg("Failed to get bookings for cancelled trip for processing refunds")
		return nil 
	}

	// 5. Refund/Cancel Bookings
	for _, b := range bookings {
		// Refund if PAID
		if b.TransactionStatus == "PAID" {
			req := &payment.RefundRequest{
				BookingID:    b.ID,
				Reason:       "Trip Cancelled by Operator",
				RefundAmount: b.TotalAmount,
			}
			_, err := s.paymentClient.CreateRefund(ctx, req)
			if err != nil {
				log.Error().Err(err).Str("booking_id", b.ID.String()).Msg("Failed to process refund for booking")
			}
		}

		// Cancel Booking
		err := s.bookingClient.CancelBooking(ctx, b.ID, "Trip Cancelled by Operator")
		if err != nil {
			log.Error().Err(err).Str("booking_id", b.ID.String()).Msg("Failed to cancel booking")
		}
	}

	return nil
}
