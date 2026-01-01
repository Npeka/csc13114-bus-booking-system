package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"bus-booking/chatbot-service/internal/model"

	"github.com/rs/zerolog/log"
)

// TripServiceClient interfaces with trip-service
type TripServiceClient interface {
	SearchTrips(ctx context.Context, params *model.TripSearchParams) (interface{}, error)
	GetTripByID(ctx context.Context, tripID string) (*model.TripDetailResponse, error)
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
	baseURL := fmt.Sprintf("%s/api/v1/trips/search", c.baseURL)

	// Build query parameters (trip-service expects GET with query params)
	queryParams := make(map[string]string)

	if params.Origin != "" {
		queryParams["origin"] = params.Origin
	}

	if params.Destination != "" {
		queryParams["destination"] = params.Destination
	}

	if !params.DepartureDate.IsZero() {
		// Format as ISO date string for the trip-service API
		queryParams["departure_time_start"] = params.DepartureDate.Format("2006-01-02T00:00:00Z")
		queryParams["departure_time_end"] = params.DepartureDate.Add(24 * time.Hour).Format("2006-01-02T00:00:00Z")
	}

	// Construct URL with query parameters
	url := baseURL
	if len(queryParams) > 0 {
		query := ""
		for key, value := range queryParams {
			if query != "" {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", key, value)
		}
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Msg("Trip service returned non-200 status")
		return nil, fmt.Errorf("trip service returned status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *tripServiceClientImpl) GetTripByID(ctx context.Context, tripID string) (*model.TripDetailResponse, error) {
	// Build URL with query parameters to preload bus, seats, and booking status
	reqURL := fmt.Sprintf("%s/api/v1/trips/%s?preload_bus=true&preload_seat=true&seat_booking_status=true&preload_route=true&preload_route_stop=true", c.baseURL, tripID)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Msg("Trip service returned non-200 status")
		return nil, fmt.Errorf("trip service returned status %d", resp.StatusCode)
	}

	var apiResp model.APIResponse[model.TripDetailResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp.Data, nil
}

// BookingServiceClient interfaces with booking-service
type BookingServiceClient interface {
	CreateGuestBooking(ctx context.Context, req *model.CreateGuestBookingRequest) (*model.BookingResponse, error)
	GetBookingByReference(ctx context.Context, reference string, email string) (*model.BookingResponse, error)
	GetBookingByID(ctx context.Context, bookingID string) (*model.BookingResponse, error)
}

type bookingServiceClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewBookingServiceClient(baseURL string) BookingServiceClient {
	return &bookingServiceClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second, // Longer timeout for booking operations
		},
	}
}

func (c *bookingServiceClientImpl) CreateGuestBooking(ctx context.Context, req *model.CreateGuestBookingRequest) (*model.BookingResponse, error) {
	url := fmt.Sprintf("%s/api/v1/bookings/guest", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call booking service")
		return nil, fmt.Errorf("failed to call booking service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Error().Int("status_code", resp.StatusCode).Msg("Booking service returned error status")

		// Try to decode error response
		var errorResp model.APIResponse[interface{}]
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Error != nil {
			return nil, fmt.Errorf("booking failed: %s", errorResp.Error.Message)
		}

		return nil, fmt.Errorf("booking service returned status %d", resp.StatusCode)
	}

	var apiResp model.APIResponse[model.BookingResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp.Data, nil
}

func (c *bookingServiceClientImpl) GetBookingByReference(ctx context.Context, reference string, email string) (*model.BookingResponse, error) {
	// Build URL with query parameters
	baseURL := fmt.Sprintf("%s/api/v1/bookings/lookup", c.baseURL)
	params := url.Values{}
	params.Add("reference", reference)
	params.Add("email", email)
	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call booking service")
		return nil, fmt.Errorf("failed to call booking service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Msg("Booking service returned non-200 status")
		return nil, fmt.Errorf("booking not found or service error: status %d", resp.StatusCode)
	}

	var apiResp model.APIResponse[model.BookingResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp.Data, nil
}

func (c *bookingServiceClientImpl) GetBookingByID(ctx context.Context, bookingID string) (*model.BookingResponse, error) {
	reqURL := fmt.Sprintf("%s/api/v1/bookings/%s", c.baseURL, bookingID)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call booking service")
		return nil, fmt.Errorf("failed to call booking service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Str("booking_id", bookingID).Msg("Booking service returned non-200 status")
		return nil, fmt.Errorf("booking not found: status %d", resp.StatusCode)
	}

	var apiResp model.APIResponse[model.BookingResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp.Data, nil
}

// PaymentServiceClient interfaces with payment-service
type PaymentServiceClient interface {
	CreateTransaction(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error)
}

type paymentServiceClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewPaymentServiceClient(baseURL string) PaymentServiceClient {
	return &paymentServiceClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *paymentServiceClientImpl) CreateTransaction(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/transactions", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call payment service")
		return nil, fmt.Errorf("failed to call payment service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Error().Int("status_code", resp.StatusCode).Msg("Payment service returned error status")

		// Try to decode error response
		var errorResp model.APIResponse[interface{}]
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Error != nil {
			return nil, fmt.Errorf("payment failed: %s", errorResp.Error.Message)
		}

		return nil, fmt.Errorf("payment service returned status %d", resp.StatusCode)
	}

	var apiResp model.APIResponse[model.TransactionResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp.Data, nil
}
