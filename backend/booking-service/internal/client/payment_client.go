package client

import (
	"bus-booking/booking-service/internal/model/payment"
	"bus-booking/shared/client"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PaymentClient interface {
	CreateTransaction(ctx context.Context, req *payment.CreateTransactionRequest) (*payment.TransactionResponse, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*payment.TransactionResponse, error)
	CancelTransaction(ctx context.Context, transactionID uuid.UUID) (*payment.TransactionResponse, error)
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

func (c *PaymentClientImpl) CreateTransaction(ctx context.Context, req *payment.CreateTransactionRequest) (*payment.TransactionResponse, error) {
	resp, err := c.http.Post(ctx, "/api/v1/transactions", req, nil)
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

func (c *PaymentClientImpl) CancelTransaction(ctx context.Context, transactionID uuid.UUID) (*payment.TransactionResponse, error) {
	resp, err := c.http.Post(ctx, fmt.Sprintf("/api/v1/transactions/%s/cancel", transactionID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment: %w", err)
	}

	transactionResp, err := client.ParseData[payment.TransactionResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cancel payment response: %w", err)
	}

	return transactionResp, nil
}
