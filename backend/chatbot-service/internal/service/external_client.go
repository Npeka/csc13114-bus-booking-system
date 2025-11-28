package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bus-booking/chatbot-service/internal/model"

	"github.com/rs/zerolog/log"
)

// TripServiceClient interfaces with trip-service
type TripServiceClient interface {
	SearchTrips(ctx context.Context, params *model.TripSearchParams) (interface{}, error)
}

type tripServiceClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewTripServiceClient(baseURL string) TripServiceClient {
	return &tripServiceClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *tripServiceClientImpl) SearchTrips(ctx context.Context, params *model.TripSearchParams) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/trips/search", c.baseURL)

	// Build query params
	reqBody := map[string]interface{}{
		"origin":      params.Origin,
		"destination": params.Destination,
	}

	if !params.DepartureDate.IsZero() {
		reqBody["departure_date"] = params.DepartureDate.Format("2006-01-02")
	}

	if params.Passengers > 0 {
		reqBody["passengers"] = params.Passengers
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call trip service")
		return nil, fmt.Errorf("failed to call trip service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// BookingServiceClient interfaces with booking-service
type BookingServiceClient interface {
	CreateBooking(ctx context.Context, bookingData interface{}) (interface{}, error)
	GetBooking(ctx context.Context, bookingID string) (interface{}, error)
}

type bookingServiceClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewBookingServiceClient(baseURL string) BookingServiceClient {
	return &bookingServiceClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *bookingServiceClientImpl) CreateBooking(ctx context.Context, bookingData interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/bookings", c.baseURL)

	jsonData, err := json.Marshal(bookingData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call booking service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *bookingServiceClientImpl) GetBooking(ctx context.Context, bookingID string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/bookings/%s", c.baseURL, bookingID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call booking service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
