package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bus-booking/shared/constants"
	sharedContext "bus-booking/shared/context"

	"github.com/rs/zerolog/log"
)

type NotificationClient interface {
	SendTripReminder(ctx context.Context, req *TripReminderRequest) error
	SendBookingConfirmation(ctx context.Context, req *BookingConfirmationRequest) error
	SendBookingFailure(ctx context.Context, req *BookingFailureRequest) error
	SendBookingPending(ctx context.Context, req *BookingPendingRequest) error
}

type TripReminderRequest struct {
	Email             string `json:"email"`
	PassengerName     string `json:"passenger_name"`
	BookingReference  string `json:"booking_reference"`
	DepartureLocation string `json:"departure_location"`
	Destination       string `json:"destination"`
	DepartureTime     string `json:"departure_time"`
	SeatNumbers       string `json:"seat_numbers"`
	BusPlate          string `json:"bus_plate"`
	PickupPoint       string `json:"pickup_point"`
	TicketLink        string `json:"ticket_link"`
}

type BookingConfirmationRequest struct {
	Email            string `json:"email"`
	Name             string `json:"name"`
	BookingReference string `json:"booking_reference"`
	From             string `json:"from"`
	To               string `json:"to"`
	DepartureTime    string `json:"departure_time"`
	SeatNumbers      string `json:"seat_numbers"`
	TotalAmount      int    `json:"total_amount"`
	TicketLink       string `json:"ticket_link"`
}

type BookingFailureRequest struct {
	Email            string `json:"email"`
	Name             string `json:"name"`
	BookingReference string `json:"booking_reference"`
	Reason           string `json:"reason"`
	From             string `json:"from"`
	To               string `json:"to"`
	DepartureTime    string `json:"departure_time"`
	BookingLink      string `json:"booking_link"`
}

type BookingPendingRequest struct {
	Email            string `json:"email"`
	Name             string `json:"name"`
	BookingReference string `json:"booking_reference"`
	From             string `json:"from"`
	To               string `json:"to"`
	DepartureTime    string `json:"departure_time"`
	TotalAmount      int    `json:"total_amount"`
	PaymentLink      string `json:"payment_link"`
}

type notificationClientImpl struct {
	baseURL     string
	serviceName string
	client      *http.Client
}

func NewNotificationClient(serviceName, baseURL string) NotificationClient {
	return &notificationClientImpl{
		baseURL:     baseURL,
		serviceName: serviceName,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *notificationClientImpl) sendRequest(ctx context.Context, path string, payload interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Propagate request ID
	reqCtx := sharedContext.FromRequestContext(ctx)
	if reqCtx.RequestID == "" {
		reqCtx.RequestID = sharedContext.GenerateRequestID()
	}

	httpReq.Header.Set(constants.XRequestID, reqCtx.RequestID)
	httpReq.Header.Set(constants.XServiceName, c.serviceName)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Error().Int("status", resp.StatusCode).Str("path", path).Msg("Notification service request failed")
		return fmt.Errorf("notification service returned status: %d", resp.StatusCode)
	}

	return nil
}

type GenericNotificationRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func (c *notificationClientImpl) toPayload(req interface{}) map[string]interface{} {
	data, err := json.Marshal(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal request to JSON")
		return nil
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal JSON to map")
		return nil
	}
	return payload
}

func (c *notificationClientImpl) SendTripReminder(ctx context.Context, req *TripReminderRequest) error {
	genReq := GenericNotificationRequest{
		Type:    "TRIP_REMINDER",
		Payload: c.toPayload(req),
	}
	return c.sendRequest(ctx, "/api/v1/notifications", genReq)
}

func (c *notificationClientImpl) SendBookingConfirmation(ctx context.Context, req *BookingConfirmationRequest) error {
	genReq := GenericNotificationRequest{
		Type:    "BOOKING_CONFIRMATION",
		Payload: c.toPayload(req),
	}
	return c.sendRequest(ctx, "/api/v1/notifications", genReq)
}

func (c *notificationClientImpl) SendBookingFailure(ctx context.Context, req *BookingFailureRequest) error {
	genReq := GenericNotificationRequest{
		Type:    "BOOKING_FAILURE",
		Payload: c.toPayload(req),
	}
	return c.sendRequest(ctx, "/api/v1/notifications", genReq)
}

func (c *notificationClientImpl) SendBookingPending(ctx context.Context, req *BookingPendingRequest) error {
	genReq := GenericNotificationRequest{
		Type:    "BOOKING_PENDING",
		Payload: c.toPayload(req),
	}
	return c.sendRequest(ctx, "/api/v1/notifications", genReq)
}
