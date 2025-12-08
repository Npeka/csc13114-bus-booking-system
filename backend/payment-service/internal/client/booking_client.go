package client

import (
	"bus-booking/payment-service/internal/model/booking"
	"bus-booking/shared/client"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type BookingClient interface {
	UpdateBookingPaymentStatus(ctx context.Context, bookingID uuid.UUID, req *booking.UpdatePaymentStatusRequest) error
}

type BookingClientImpl struct {
	http client.HTTPClient
}

func NewBookingClient(serviceName, baseURL string) BookingClient {
	httpClient := client.NewHTTPClient(&client.Config{
		ServiceName: serviceName,
		BaseURL:     baseURL,
	})

	return &BookingClientImpl{
		http: httpClient,
	}
}

func (c *BookingClientImpl) UpdateBookingPaymentStatus(ctx context.Context, bookingID uuid.UUID, req *booking.UpdatePaymentStatusRequest) error {
	endpoint := fmt.Sprintf("/api/v1/bookings/%s/payment-status", bookingID.String())
	_, err := c.http.Put(ctx, endpoint, req, nil)
	if err != nil {
		return fmt.Errorf("failed to update booking payment status: %w", err)
	}
	return nil
}
