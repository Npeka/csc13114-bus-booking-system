package client

import (
	"context"
	"fmt"

	"bus-booking/shared/client"
	"bus-booking/trip-service/internal/model/payment"
)

type PaymentClient interface {
	CreateRefund(ctx context.Context, req *payment.RefundRequest) (*payment.RefundResponse, error)
}

type paymentClientImpl struct {
	httpClient client.HTTPClient
	baseURL    string
}

func NewPaymentClient(serviceName, baseURL string) PaymentClient {
	return &paymentClientImpl{
		httpClient: client.NewHTTPClient(&client.Config{
			ServiceName: serviceName,
			BaseURL:     baseURL,
		}),
		baseURL: baseURL,
	}
}

func (c *paymentClientImpl) CreateRefund(ctx context.Context, req *payment.RefundRequest) (*payment.RefundResponse, error) {
	url := "/api/v1/refunds"
	resp, err := c.httpClient.Post(ctx, url, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	refund, err := client.ParseData[payment.RefundResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refund response: %w", err)
	}

	return refund, nil
}
