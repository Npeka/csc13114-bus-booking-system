package service

import (
	"context"
	"fmt"
	"time"

	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/model"

	"github.com/payOSHQ/payos-lib-golang/v2"
)

type PayOSService interface {
	CreatePaymentLink(ctx context.Context, req *model.CreatePayOSPaymentLinkRequest) (*payos.CreatePaymentLinkResponse, error)
	GetPaymentLink(ctx context.Context, paymentLinkID string) (*payos.PaymentLink, error)
	VerifyWebhook(ctx context.Context, webhookData map[string]interface{}) error
	CancelPaymentLink(ctx context.Context, paymentLinkID string, cancellationReason *string) (*payos.PaymentLink, error)
	ToTransactionStatus(payOSStatus payos.PaymentLinkStatus) model.TransactionStatus
}

// PayOSServiceImpl handles PayOS API integration
type PayOSServiceImpl struct {
	payOSClient *payos.PayOS
	ReturnURL   string
	CancelURL   string
}

// NewPayOSService creates a new PayOS service
func NewPayOSService(cfg config.PayOSConfig) PayOSService {
	payOSClient, err := payos.NewPayOS(&payos.PayOSOptions{
		ClientId:    cfg.ClientID,
		ApiKey:      cfg.APIKey,
		ChecksumKey: cfg.ChecksumKey,
	})
	if err != nil {
		panic(err)
	}
	return &PayOSServiceImpl{
		payOSClient: payOSClient,
		ReturnURL:   cfg.ReturnURL,
		CancelURL:   cfg.CancelURL,
	}
}

func (c *PayOSServiceImpl) CreatePaymentLink(ctx context.Context, req *model.CreatePayOSPaymentLinkRequest) (*payos.CreatePaymentLinkResponse, error) {
	paymentLinkRequest := payos.CreatePaymentLinkRequest{
		OrderCode:   c.generateOrderCode(),
		Amount:      req.Amount,
		Description: req.Description,
		CancelUrl:   c.CancelURL,
		ReturnUrl:   c.ReturnURL,
	}

	paymentLinkResponse, err := c.payOSClient.PaymentRequests.Create(ctx, paymentLinkRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	return paymentLinkResponse, nil
}

func (c *PayOSServiceImpl) GetPaymentLink(ctx context.Context, paymentLinkID string) (*payos.PaymentLink, error) {
	return c.payOSClient.PaymentRequests.Get(ctx, paymentLinkID)
}

func (c *PayOSServiceImpl) generateOrderCode() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (c *PayOSServiceImpl) VerifyWebhook(ctx context.Context, webhookData map[string]interface{}) error {
	_, err := c.payOSClient.Webhooks.VerifyData(ctx, webhookData)
	if err != nil {
		return err
	}
	return nil
}

func (c *PayOSServiceImpl) CancelPaymentLink(ctx context.Context, paymentLinkID string, cancellationReason *string) (*payos.PaymentLink, error) {
	return c.payOSClient.PaymentRequests.Cancel(ctx, paymentLinkID, cancellationReason)
}

func (c *PayOSServiceImpl) ToTransactionStatus(payOSStatus payos.PaymentLinkStatus) model.TransactionStatus {
	switch payOSStatus {
	case payos.PaymentLinkStatusPending:
		return model.TransactionStatusPending
	case payos.PaymentLinkStatusCancelled:
		return model.TransactionStatusCancelled
	case payos.PaymentLinkStatusUnderpaid:
		return model.TransactionStatusUnderpaid
	case payos.PaymentLinkStatusPaid:
		return model.TransactionStatusPaid
	case payos.PaymentLinkStatusExpired:
		return model.TransactionStatusExpired
	case payos.PaymentLinkStatusProcessing:
		return model.TransactionStatusProcessing
	case payos.PaymentLinkStatusFailed:
		return model.TransactionStatusFailed
	default:
		return model.TransactionStatusPending
	}
}
