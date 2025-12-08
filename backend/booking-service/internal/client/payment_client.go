package client

import (
	"bus-booking/booking-service/internal/model/payment"
	"bus-booking/shared/client"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PaymentClient interface {
	CreatePaymentLink(ctx context.Context, req *payment.CreatePaymentLinkRequest) (*payment.TransactionResponse, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*payment.TransactionResponse, error)
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

func (c *PaymentClientImpl) CreatePaymentLink(ctx context.Context, req *payment.CreatePaymentLinkRequest) (*payment.TransactionResponse, error) {
	resp, err := c.http.Post(ctx, "/api/v1/transactions/payment-link", req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	paymentResp, err := client.ParseData[payment.TransactionResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment response: %w", err)
	}

	return paymentResp, nil
}

func (c *PaymentClientImpl) GetTransactionByID(ctx context.Context, id uuid.UUID) (*payment.TransactionResponse, error) {
	resp, err := c.http.Get(ctx, fmt.Sprintf("/api/v1/transactions/%s", id), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	transactionResp, err := client.ParseData[payment.TransactionResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction response: %w", err)
	}

	return transactionResp, nil
}
