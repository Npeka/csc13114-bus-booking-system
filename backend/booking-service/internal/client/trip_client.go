package client

import (
	"context"
	"fmt"

	"bus-booking/booking-service/internal/model/trip"
	"bus-booking/shared/client"

	"github.com/google/uuid"
)

type TripClient interface {
	GetTripByID(ctx context.Context, req trip.GetTripByIDRequest, ripID uuid.UUID) (*trip.Trip, error)
	GetTripsByIDs(ctx context.Context, req trip.GetTripByIDRequest, tripIDs []uuid.UUID) ([]trip.Trip, error)
	ListSeatsByIDs(ctx context.Context, seatIDs []uuid.UUID) ([]trip.Seat, error)
}

type TripClientImpl struct {
	http client.HTTPClient
}

func NewTripClient(serviceName, baseURL string) TripClient {
	httpClient := client.NewHTTPClient(&client.Config{
		ServiceName: serviceName,
		BaseURL:     baseURL,
	})

	return &TripClientImpl{
		http: httpClient,
	}
}

func (c *TripClientImpl) GetTripByID(ctx context.Context, req trip.GetTripByIDRequest, tripID uuid.UUID) (*trip.Trip, error) {
	endpoint := fmt.Sprintf("/api/v1/trips/%s", tripID.String())

	params := make(map[string][]string)
	if req.SeatBookingStatus {
		params["seat_booking_status"] = []string{"true"}
	}
	if req.PreLoadRoute {
		params["preload_route"] = []string{"true"}
	}
	if req.PreLoadRouteStop {
		params["preload_route_stop"] = []string{"true"}
	}
	if req.PreloadBus {
		params["preload_bus"] = []string{"true"}
	}
	if req.PreloadSeat {
		params["preload_seat"] = []string{"true"}
	}

	res, err := c.http.Get(ctx, endpoint, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	tripData, err := client.ParseData[trip.Trip](res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trip response: %w", err)
	}

	return tripData, nil
}

func (c *TripClientImpl) ListSeatsByIDs(ctx context.Context, seatIDs []uuid.UUID) ([]trip.Seat, error) {
	endpoint := "/api/v1/buses/seats/ids"

	params := make(map[string][]string)
	seatIDStrings := make([]string, len(seatIDs))
	for i, id := range seatIDs {
		seatIDStrings[i] = id.String()
	}
	params["seat_ids"] = seatIDStrings

	res, err := c.http.Get(ctx, endpoint, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	seats, err := client.ParseListData[trip.Seat](res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trip response: %w", err)
	}

	return seats, nil
}

// GetTripsByIDs fetches multiple trips by IDs in a single batch request
func (c *TripClientImpl) GetTripsByIDs(ctx context.Context, req trip.GetTripByIDRequest, tripIDs []uuid.UUID) ([]trip.Trip, error) {
	endpoint := "/api/v1/trips"

	params := make(map[string][]string)
	// Add trip IDs
	tripIDStrings := make([]string, len(tripIDs))
	for i, id := range tripIDs {
		tripIDStrings[i] = id.String()
	}
	params["ids[]"] = tripIDStrings

	// Add preload options
	if req.PreLoadRoute {
		params["preload_route"] = []string{"true"}
	}
	if req.PreloadBus {
		params["preload_bus"] = []string{"true"}
	}

	res, err := c.http.Get(ctx, endpoint, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips: %w", err)
	}

	trips, err := client.ParseListData[trip.Trip](res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trips response: %w", err)
	}

	return trips, nil
}
