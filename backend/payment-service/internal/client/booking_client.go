package client

import (
	"bus-booking/shared/client"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UpdatePaymentStatusRequest struct {
	PaymentStatus  string `json:"payment_status"`
	BookingStatus  string `json:"booking_status"`
	PaymentOrderID string `json:"payment_order_id"`
}

type BookingClient interface {
	UpdateBookingPaymentStatus(ctx context.Context, bookingID uuid.UUID, req *UpdatePaymentStatusRequest) error
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

func (c *BookingClientImpl) UpdateBookingPaymentStatus(ctx context.Context, bookingID uuid.UUID, req *UpdatePaymentStatusRequest) error {
	endpoint := fmt.Sprintf("/api/v1/bookings/%s/payment-status", bookingID.String())
	_, err := c.http.Put(ctx, endpoint, req, nil)
	if err != nil {
		return fmt.Errorf("failed to update booking payment status: %w", err)
	}
	return nil
}
