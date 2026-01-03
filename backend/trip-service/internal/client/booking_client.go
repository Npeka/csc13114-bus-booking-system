package client

import (
	"context"
	"fmt"

	"bus-booking/shared/client"
	"bus-booking/trip-service/internal/model/booking"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BookingClient interface {
	GetSeatStatus(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) ([]booking.SeatStatus, error)
	GetTripBookings(ctx context.Context, tripID uuid.UUID) ([]*booking.Booking, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID, reason string) error
}

type bookingClientImpl struct {
	httpClient client.HTTPClient
	baseURL    string
}

func NewBookingClient(serviceName, baseURL string) BookingClient {
	return &bookingClientImpl{
		httpClient: client.NewHTTPClient(&client.Config{
			ServiceName: serviceName,
			BaseURL:     baseURL,
		}),
		baseURL: baseURL,
	}
}

func (c *bookingClientImpl) GetSeatStatus(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) ([]booking.SeatStatus, error) {
	if len(seatIDs) == 0 {
		return []booking.SeatStatus{}, nil
	}

	// Build query params
	params := make(map[string][]string)
	seatIDStrings := make([]string, len(seatIDs))
	for i, id := range seatIDs {
		seatIDStrings[i] = id.String()
	}
	params["seat_ids"] = seatIDStrings

	url := fmt.Sprintf("/api/v1/bookings/trips/%s/seats/status", tripID)

	resp, err := c.httpClient.Get(ctx, url, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat status: %w", err)
	}

	seatStatuses, err := client.ParseListData[booking.SeatStatus](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse seat status response: %w", err)
	}
	log.Info().Msgf("Received seat statuses: %+v", seatStatuses)
	return seatStatuses, nil
}

func (c *bookingClientImpl) GetTripBookings(ctx context.Context, tripID uuid.UUID) ([]*booking.Booking, error) {
	// Fetch all bookings (using larger limit)
	params := map[string][]string{
		"page":  {"1"},
		"limit": {"1000"},
	}

	url := fmt.Sprintf("/api/v1/bookings/trip/%s", tripID)
	resp, err := c.httpClient.Get(ctx, url, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip bookings: %w", err)
	}

	// Parse paginated response
	paginatedResp, err := client.ParseData[booking.BookingListResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse booking list: %w", err)
	}

	return paginatedResp.Data, nil
}

func (c *bookingClientImpl) CancelBooking(ctx context.Context, bookingID uuid.UUID, reason string) error {
	req := booking.CancelBookingRequest{
		Reason: reason,
	}

	url := fmt.Sprintf("/api/v1/bookings/%s/cancel", bookingID)
	_, err := c.httpClient.Post(ctx, url, req, nil)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	return nil
}
