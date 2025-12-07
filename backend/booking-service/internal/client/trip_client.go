package client

import (
	"context"
	"fmt"

	"bus-booking/booking-service/internal/model/trip"
	"bus-booking/shared/client"

	"github.com/google/uuid"
)

type TripClient interface {
	GetTrip(ctx context.Context, tripID uuid.UUID) (*trip.Trip, error)
	GetSeat(ctx context.Context, tripID, seatID uuid.UUID) (*trip.Seat, error)
	GetSeatsMetadata(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) (map[uuid.UUID]*trip.Seat, error)
	CalculateTotalPrice(basePrice float64, seats map[uuid.UUID]*trip.Seat) float64
}

type TripClientImpl struct {
	http client.HTTPClient
}

func NewTripClient(ServiceName, baseURL string) TripClient {
	httpClient := client.NewHTTPClient(&client.Config{
		ServiceName: ServiceName,
		BaseURL:     baseURL,
	})

	return &TripClientImpl{
		http: httpClient,
	}
}

func (c *TripClientImpl) GetTrip(ctx context.Context, tripID uuid.UUID) (*trip.Trip, error) {
	endpoint := fmt.Sprintf("/api/v1/trips/%s", tripID.String())

	resp, err := c.http.Get(ctx, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	tripData, err := client.ParseData[trip.Trip](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trip response: %w", err)
	}

	return tripData, nil
}

// GetSeat finds specific seat by ID from trip's bus seats
func (c *TripClientImpl) GetSeat(ctx context.Context, tripID, seatID uuid.UUID) (*trip.Seat, error) {
	tripData, err := c.GetTrip(ctx, tripID)
	if err != nil {
		return nil, err
	}

	if tripData.Bus == nil {
		return nil, fmt.Errorf("bus information not available")
	}

	for i := range tripData.Bus.Seats {
		if tripData.Bus.Seats[i].ID == seatID {
			return &tripData.Bus.Seats[i], nil
		}
	}

	return nil, fmt.Errorf("seat %s not found", seatID)
}

// GetSeatsMetadata retrieves seat metadata (number, type, price multiplier) from trip service
// Does NOT check booking status - that's handled by booking service
func (c *TripClientImpl) GetSeatsMetadata(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) (map[uuid.UUID]*trip.Seat, error) {
	tripData, err := c.GetTrip(ctx, tripID)
	if err != nil {
		return nil, err
	}

	if !tripData.IsBookable() {
		return nil, fmt.Errorf("trip is not available for booking")
	}

	if tripData.Bus == nil || len(tripData.Bus.Seats) == 0 {
		return nil, fmt.Errorf("bus seats not available")
	}

	// Create seat map
	seatMap := make(map[uuid.UUID]*trip.Seat)
	for i := range tripData.Bus.Seats {
		seatMap[tripData.Bus.Seats[i].ID] = &tripData.Bus.Seats[i]
	}

	// Get metadata for all requested seats
	validatedSeats := make(map[uuid.UUID]*trip.Seat)
	for _, seatID := range seatIDs {
		seat, exists := seatMap[seatID]
		if !exists {
			return nil, fmt.Errorf("seat %s not found", seatID)
		}
		// Only check if seat exists in bus, not booking status
		validatedSeats[seatID] = seat
	}

	return validatedSeats, nil
}

// CalculateTotalPrice calculates total price for seats
func (c *TripClientImpl) CalculateTotalPrice(basePrice float64, seats map[uuid.UUID]*trip.Seat) float64 {
	total := 0.0
	for _, seat := range seats {
		total += seat.CalculateSeatPrice(basePrice)
	}
	return total
}
