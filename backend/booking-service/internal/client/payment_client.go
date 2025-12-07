package client

import (
	"bus-booking/shared/client"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreatePaymentLinkRequest struct {
	BookingID     uuid.UUID `json:"booking_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Description   string    `json:"description"`
	BuyerName     string    `json:"buyer_name"`
	BuyerEmail    string    `json:"buyer_email"`
	BuyerPhone    string    `json:"buyer_phone"`
}

type PaymentLinkResponse struct {
	ID          uuid.UUID `json:"id"`
	BookingID   uuid.UUID `json:"booking_id"`
	OrderCode   int64     `json:"order_code"`
	Status      string    `json:"status"`
	CheckoutURL string    `json:"checkout_url"`
	QRCode      string    `json:"qr_code"`
}

type PaymentClient interface {
	CreatePaymentLink(ctx context.Context, req *CreatePaymentLinkRequest) (*PaymentLinkResponse, error)
}

type PaymentClientImpl struct {
	http client.HTTPClient
}

func NewPaymentClient(serviceName, baseURL string) PaymentClient {
	httpClient := client.NewHTTPClient(&client.Config{
		ServiceName: serviceName,
		BaseURL:     baseURL,
	})

	return &PaymentClientImpl{
		http: httpClient,
	}
}

func (c *PaymentClientImpl) CreatePaymentLink(ctx context.Context, req *CreatePaymentLinkRequest) (*PaymentLinkResponse, error) {
	resp, err := c.http.Post(ctx, "/api/v1/transactions/payment-link", req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	paymentResp, err := client.ParseData[PaymentLinkResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment response: %w", err)
	}

	return paymentResp, nil
}
