package service

import (
	"bus-booking/payment-service/internal/client"
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/model/booking"
	"bus-booking/payment-service/internal/repository"
	"bus-booking/shared/ginext"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/payOSHQ/payos-lib-golang/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type TransactionService interface {
	CreatePaymentLink(ctx context.Context, req *model.CreatePaymentLinkRequest, userID uuid.UUID) (*model.TransactionResponse, error)
	HandlePaymentWebhook(ctx context.Context, webhookData map[string]interface{}) error
	CancelPayment(ctx context.Context, transactionID uuid.UUID) (*model.TransactionResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TransactionResponse, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.TransactionResponse, error)
}

type TransactionServiceImpl struct {
	transactionRepo repository.TransactionRepository
	bookingClient   client.BookingClient
	payOSService    PayOSService
}

func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	bookingClient client.BookingClient,
	payOSService PayOSService,
) TransactionService {
	return &TransactionServiceImpl{
		transactionRepo: transactionRepo,
		bookingClient:   bookingClient,
		payOSService:    payOSService,
	}
}

func (s *TransactionServiceImpl) CreatePaymentLink(ctx context.Context, req *model.CreatePaymentLinkRequest, userID uuid.UUID) (*model.TransactionResponse, error) {
	payosResp, err := s.payOSService.CreatePaymentLink(ctx, &model.CreatePayOSPaymentLinkRequest{
		Amount:      req.Amount,
		Description: req.Description,
	})
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to create payment link: %v", err))
	}

	transaction := &model.Transaction{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		BookingID:     req.BookingID,
		UserID:        userID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		OrderCode:     payosResp.OrderCode,
		PaymentLinkID: payosResp.PaymentLinkId,
		Status:        s.payOSService.ToTransactionStatus(payosResp.Status),
		CheckoutURL:   payosResp.CheckoutUrl,
		QRCode:        payosResp.QrCode,
	}

	if err = s.transactionRepo.CreateTransaction(ctx, transaction); err != nil {
		log.Error().Err(err).Msg("Failed to save transaction")
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	return s.toTransactionResponse(transaction), nil
}

func (s *TransactionServiceImpl) HandlePaymentWebhook(ctx context.Context, webhookData map[string]interface{}) error {
	log.Info().Msg("Starting webhook verification")
	if err := s.payOSService.VerifyWebhook(ctx, webhookData); err != nil {
		log.Error().Err(err).Msg("Webhook signature verification failed")
		return ginext.NewUnauthorizedError("invalid webhook signature")
	}
	log.Info().Msg("Webhook signature verified successfully")

	// Convert map to struct: marshal to JSON bytes, then unmarshal to struct
	log.Info().Msg("Converting webhook data to struct")
	webhookBytes, err := json.Marshal(webhookData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal webhook data to JSON bytes")
		return ginext.NewBadRequestError("invalid webhook data format")
	}
	log.Info().Str("webhookJSON", string(webhookBytes)).Msg("Webhook data marshaled successfully")

	var webhookDataModel model.PaymentWebhookData
	if err := json.Unmarshal(webhookBytes, &webhookDataModel); err != nil {
		log.Error().Err(err).Str("webhookJSON", string(webhookBytes)).Msg("Failed to unmarshal webhook data to struct")
		return ginext.NewBadRequestError("invalid webhook data")
	}
	log.Info().Interface("webhookDataModel", webhookDataModel).Msg("Webhook data unmarshaled successfully")

	var (
		paymentLink *payos.PaymentLink
		transaction *model.Transaction
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		paymentLink, err = s.payOSService.GetPaymentLink(gCtx, webhookDataModel.Data.PaymentLinkID)
		return err
	})

	g.Go(func() error {
		var err error
		orderCode := webhookDataModel.Data.OrderCode
		paymentLinkID := webhookDataModel.Data.PaymentLinkID
		transaction, err = s.transactionRepo.GetByWebhookData(ctx, orderCode, paymentLinkID)
		return err
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// Update transaction status
	transaction.Status = s.payOSService.ToTransactionStatus(paymentLink.Status)
	transaction.Reference = webhookDataModel.Data.Reference

	// Parse transaction datetime
	if transTime, err := time.Parse("2006-01-02 15:04:05", webhookDataModel.Data.TransactionDateTime); err == nil {
		transTimeUnix := transTime.Unix()
		transaction.TransactionTime = &transTimeUnix
	}

	// Update in database
	if err := s.transactionRepo.UpdateTransaction(ctx, transaction); err != nil {
		log.Error().Err(err).Msg("Failed to update transaction")
		return ginext.NewInternalServerError("failed to update transaction")
	}

	// Notify booking service about payment success
	if err := s.bookingClient.UpdateBookingStatus(ctx, &booking.UpdateBookingStatusRequest{
		TransactionStatus: transaction.Status,
	}, transaction.BookingID); err != nil {
		log.Error().Err(err).
			Str("booking_id", transaction.BookingID.String()).
			Msg("Failed to update booking payment status")
		// Don't fail the webhook - payment is already recorded
		// This can be retried via a background job
	}

	return nil
}
func (s *TransactionServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ginext.NewNotFoundError("transaction not found")
	}
	return s.toTransactionResponse(transaction), nil
}

func (s *TransactionServiceImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByBookingID(ctx, bookingID)
	if err != nil {
		return nil, ginext.NewNotFoundError("transaction not found")
	}
	return s.toTransactionResponse(transaction), nil
}

// CancelPayment cancels a payment and updates transaction status
func (s *TransactionServiceImpl) CancelPayment(ctx context.Context, transactionID uuid.UUID) (*model.TransactionResponse, error) {
	// Get transaction from DB
	transaction, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, ginext.NewNotFoundError("transaction not found")
	}

	// Check if already cancelled or completed
	if transaction.Status == model.TransactionStatusCancelled {
		return nil, ginext.NewBadRequestError("transaction already cancelled")
	}

	if transaction.Status == model.TransactionStatusPaid {
		return nil, ginext.NewBadRequestError("cannot cancel a paid transaction")
	}

	// Cancel PayOS payment link
	reason := "Booking cancelled by user"
	paymentLink, err := s.payOSService.CancelPaymentLink(ctx, transaction.PaymentLinkID, &reason)
	if err != nil {
		log.Error().Err(err).
			Str("transaction_id", transactionID.String()).
			Str("payment_link_id", transaction.PaymentLinkID).
			Msg("Failed to cancel PayOS payment link")
		// Continue to update local status even if PayOS call fails
	}

	// Update transaction status
	if paymentLink != nil {
		transaction.Status = s.payOSService.ToTransactionStatus(paymentLink.Status)
	} else {
		// If PayOS call failed, mark as cancelled locally
		transaction.Status = model.TransactionStatusCancelled
	}

	if err := s.transactionRepo.UpdateTransaction(ctx, transaction); err != nil {
		log.Error().Err(err).Msg("Failed to update transaction status")
		return nil, ginext.NewInternalServerError("failed to update transaction")
	}

	log.Info().
		Str("transaction_id", transactionID.String()).
		Str("new_status", string(transaction.Status)).
		Msg("Successfully cancelled payment")

	return s.toTransactionResponse(transaction), nil
}

func (s *TransactionServiceImpl) toTransactionResponse(t *model.Transaction) *model.TransactionResponse {
	return &model.TransactionResponse{
		ID:            t.ID,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		BookingID:     t.BookingID,
		UserID:        t.UserID,
		Amount:        t.Amount,
		Currency:      t.Currency,
		PaymentMethod: t.PaymentMethod,
		OrderCode:     t.OrderCode,
		Status:        t.Status,
		CheckoutURL:   t.CheckoutURL,
		QRCode:        t.QRCode,
	}
}
