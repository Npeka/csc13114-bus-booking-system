package client

import (
	"context"
	"fmt"

	"bus-booking/booking-service/internal/model/trip"
	"bus-booking/shared/client"

	"github.com/google/uuid"
)

type TripClient interface {
	GetTripByID(ctx context.Context, tripID uuid.UUID) (*trip.Trip, error)
	ListSeatsByIDs(ctx context.Context, seatIDs []uuid.UUID) ([]trip.Seat, error)
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

func (c *TripClientImpl) GetTripByID(ctx context.Context, tripID uuid.UUID) (*trip.Trip, error) {
	endpoint := fmt.Sprintf("/api/v1/trips/%s", tripID.String())

	resp, err := c.http.Get(ctx, endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	tripData, err := client.ParseData[trip.Trip](resp)
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

	resp, err := c.http.Get(ctx, endpoint, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	seats, err := client.ParseListData[trip.Seat](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trip response: %w", err)
	}

	return seats, nil
}
