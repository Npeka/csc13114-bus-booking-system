package service

import (
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error)
	CreatePaymentLink(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error)
	GetTransactionByOrderCode(ctx context.Context, orderCode int64) (*model.Transaction, error)
	GetTransactionByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error)
	HandlePaymentWebhook(ctx context.Context, webhookData *model.PaymentWebhookData) error
	ConfirmPayment(ctx context.Context, orderCode int64) error
	CancelPayment(ctx context.Context, orderCode int64, reason string) error
}

type TransactionServiceImpl struct {
	repositories *repository.Repositories
	payosClient  *PayOSClient
	returnURL    string
	cancelURL    string
}

func NewTransactionService(repositories *repository.Repositories, payosClient *PayOSClient, returnURL, cancelURL string) TransactionService {
	return &TransactionServiceImpl{
		repositories: repositories,
		payosClient:  payosClient,
		returnURL:    returnURL,
		cancelURL:    cancelURL,
	}
}

// CreateTransaction creates a basic transaction record (deprecated - use CreatePaymentLink)
func (s *TransactionServiceImpl) CreateTransaction(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	transaction := &model.Transaction{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		BookingID:     req.BookingID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Status:        model.PaymentStatusPending,
	}

	err := s.repositories.Transaction.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create transaction")
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return s.toTransactionResponse(transaction), nil
}

// CreatePaymentLink creates a transaction and PayOS payment link
func (s *TransactionServiceImpl) CreatePaymentLink(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	// Generate unique order code from timestamp
	orderCode := time.Now().Unix()

	// Create PayOS payment link request
	payosReq := &model.CreatePaymentLinkRequest{
		OrderCode:   orderCode,
		Amount:      int(req.Amount), // Convert to int (VND doesn't have decimals)
		Description: req.Description,
		BuyerName:   req.BuyerName,
		BuyerEmail:  req.BuyerEmail,
		BuyerPhone:  req.BuyerPhone,
		CancelURL:   s.cancelURL,
		ReturnURL:   s.returnURL,
		ExpiredAt:   time.Now().Add(15 * time.Minute).Unix(), // 15 minutes expiry
	}

	// Add default description if empty
	if payosReq.Description == "" {
		payosReq.Description = fmt.Sprintf("Payment for booking %s", req.BookingID.String())
	}

	// Call PayOS API
	payosResp, err := s.payosClient.CreatePaymentLink(payosReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create PayOS payment link")
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	// Check PayOS response
	if payosResp.Code != model.PayOSCodeSuccess {
		return nil, fmt.Errorf("PayOS error: %s - %s", payosResp.Code, payosResp.Desc)
	}

	// Create transaction record
	transaction := &model.Transaction{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		BookingID:     req.BookingID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		OrderCode:     orderCode,
		PaymentLinkID: payosResp.Data.PaymentLinkID,
		Status:        payosResp.Data.Status,
		CheckoutURL:   payosResp.Data.CheckoutURL,
		QRCode:        payosResp.Data.QRCode,
	}

	// Save to database
	err = s.repositories.Transaction.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save transaction")
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	log.Info().
		Int64("order_code", orderCode).
		Str("booking_id", req.BookingID.String()).
		Str("payment_link_id", payosResp.Data.PaymentLinkID).
		Msg("Payment link created successfully")

	return s.toTransactionResponse(transaction), nil
}

// GetTransactionByOrderCode retrieves transaction by PayOS order code
func (s *TransactionServiceImpl) GetTransactionByOrderCode(ctx context.Context, orderCode int64) (*model.Transaction, error) {
	transaction, err := s.repositories.Transaction.GetTransactionByOrderCode(ctx, orderCode)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return transaction, nil
}

// GetTransactionByBookingID retrieves transaction by booking ID
func (s *TransactionServiceImpl) GetTransactionByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error) {
	transaction, err := s.repositories.Transaction.GetTransactionByBookingID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return transaction, nil
}

// HandlePaymentWebhook processes PayOS webhook notification
func (s *TransactionServiceImpl) HandlePaymentWebhook(ctx context.Context, webhookData *model.PaymentWebhookData) error {
	// Verify webhook signature
	if !s.payosClient.VerifyWebhookSignature(webhookData) {
		log.Error().Msg("Invalid webhook signature")
		return fmt.Errorf("invalid webhook signature")
	}

	// Get transaction by order code
	transaction, err := s.GetTransactionByOrderCode(ctx, webhookData.Data.OrderCode)
	if err != nil {
		log.Error().Err(err).Int64("order_code", webhookData.Data.OrderCode).Msg("Transaction not found")
		return err
	}

	// Update transaction status
	transaction.Status = model.PaymentStatusPaid
	transaction.Reference = webhookData.Data.Reference

	// Parse transaction datetime
	if transTime, err := time.Parse("2006-01-02 15:04:05", webhookData.Data.TransactionDateTime); err == nil {
		transTimeUnix := transTime.Unix()
		transaction.TransactionTime = &transTimeUnix
	}

	// Update in database
	err = s.repositories.Transaction.UpdateTransaction(ctx, transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update transaction")
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// TODO: Notify booking service about payment success
	// This should be done via message queue or HTTP call to booking service
	log.Info().
		Str("booking_id", transaction.BookingID.String()).
		Int64("order_code", transaction.OrderCode).
		Str("reference", transaction.Reference).
		Msg("Payment confirmed via webhook")

	return nil
}

// ConfirmPayment confirms payment by checking with PayOS
func (s *TransactionServiceImpl) ConfirmPayment(ctx context.Context, orderCode int64) error {
	// Get payment info from PayOS
	paymentInfo, err := s.payosClient.GetPaymentInfo(orderCode)
	if err != nil {
		return fmt.Errorf("failed to get payment info: %w", err)
	}

	// Get transaction from database
	transaction, err := s.GetTransactionByOrderCode(ctx, orderCode)
	if err != nil {
		return err
	}

	// Update transaction status based on PayOS response
	if paymentInfo.Code == model.PayOSCodeSuccess {
		transaction.Status = paymentInfo.Data.Status

		// Update reference if payment is completed
		if len(paymentInfo.Data.Transactions) > 0 {
			transaction.Reference = paymentInfo.Data.Transactions[0].Reference
			transTime := paymentInfo.Data.Transactions[0].TransactionDateTime.Unix()
			transaction.TransactionTime = &transTime
		}

		err = s.repositories.Transaction.UpdateTransaction(ctx, transaction)
		if err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		log.Info().
			Int64("order_code", orderCode).
			Str("status", transaction.Status).
			Msg("Payment status confirmed")
	}

	return nil
}

// CancelPayment cancels a payment
func (s *TransactionServiceImpl) CancelPayment(ctx context.Context, orderCode int64, reason string) error {
	// Cancel via PayOS API
	_, err := s.payosClient.CancelPayment(orderCode, reason)
	if err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	// Update transaction status
	transaction, err := s.GetTransactionByOrderCode(ctx, orderCode)
	if err != nil {
		return err
	}

	transaction.Status = model.PaymentStatusCancelled
	err = s.repositories.Transaction.UpdateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	log.Info().
		Int64("order_code", orderCode).
		Str("reason", reason).
		Msg("Payment cancelled")

	return nil
}

// toTransactionResponse converts Transaction model to response
func (s *TransactionServiceImpl) toTransactionResponse(t *model.Transaction) *model.TransactionResponse {
	return &model.TransactionResponse{
		ID:            t.ID,
		BookingID:     t.BookingID,
		Amount:        t.Amount,
		Currency:      t.Currency,
		PaymentMethod: t.PaymentMethod,
		OrderCode:     t.OrderCode,
		Status:        t.Status,
		CheckoutURL:   t.CheckoutURL,
		QRCode:        t.QRCode,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}
