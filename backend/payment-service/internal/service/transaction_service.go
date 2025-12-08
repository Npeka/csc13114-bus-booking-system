package service

import (
	"bus-booking/payment-service/internal/client"
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/model/booking"
	"bus-booking/payment-service/internal/repository"
	"bus-booking/shared/ginext"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/payOSHQ/payos-lib-golang/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type TransactionService interface {
	CreatePaymentLink(ctx context.Context, req *model.CreatePaymentLinkRequest, userID uuid.UUID) (*model.TransactionResponse, error)
	HandlePaymentWebhook(ctx context.Context, webhookData *model.PaymentWebhookData) error
	GetTransactionByOrderCode(ctx context.Context, orderCode int) (*model.Transaction, error)
	GetTransactionByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
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

func (s *TransactionServiceImpl) HandlePaymentWebhook(ctx context.Context, webhookData *model.PaymentWebhookData) error {
	if err := s.payOSService.VerifyWebhook(ctx, webhookData); err != nil {
		return ginext.NewUnauthorizedError("invalid webhook signature")
	}

	var (
		paymentLink *payos.PaymentLink
		transaction *model.Transaction
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		paymentLink, err = s.payOSService.GetPaymentLink(gCtx, webhookData.Data.PaymentLinkID)
		return err
	})

	g.Go(func() error {
		var err error
		transaction, err = s.GetTransactionByOrderCode(gCtx, webhookData.Data.OrderCode)
		return err
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// Update transaction status
	transaction.Status = s.payOSService.ToTransactionStatus(paymentLink.Status)
	transaction.Reference = webhookData.Data.Reference

	// Parse transaction datetime
	if transTime, err := time.Parse("2006-01-02 15:04:05", webhookData.Data.TransactionDateTime); err == nil {
		transTimeUnix := transTime.Unix()
		transaction.TransactionTime = &transTimeUnix
	}

	// Update in database
	if err := s.transactionRepo.UpdateTransaction(ctx, transaction); err != nil {
		log.Error().Err(err).Msg("Failed to update transaction")
		return ginext.NewInternalServerError("failed to update transaction")
	}

	// Notify booking service about payment success
	updateReq := &booking.UpdatePaymentStatusRequest{
		PaymentStatus:  transaction.Status,
		BookingStatus:  s.toBookingStatus(transaction.Status),
		TransactionID:  transaction.ID,
		PaymentOrderID: fmt.Sprintf("%d", transaction.OrderCode),
	}

	if err := s.bookingClient.UpdateBookingPaymentStatus(ctx, transaction.BookingID, updateReq); err != nil {
		log.Error().Err(err).
			Str("booking_id", transaction.BookingID.String()).
			Msg("Failed to update booking payment status")
		// Don't fail the webhook - payment is already recorded
		// This can be retried via a background job
	} else {
		log.Info().
			Str("booking_id", transaction.BookingID.String()).
			Int64("order_code", transaction.OrderCode).
			Str("reference", transaction.Reference).
			Msg("Payment confirmed and booking updated")
	}

	return nil
}

func (s *TransactionServiceImpl) toBookingStatus(status model.TransactionStatus) booking.BookingStatus {
	switch status {
	case model.TransactionStatusCancelled:
		return booking.BookingStatusCancelled
	case model.TransactionStatusUnderpaid:
		return booking.BookingStatusPending
	case model.TransactionStatusPaid:
		return booking.BookingStatusConfirmed
	case model.TransactionStatusFailed:
		return booking.BookingStatusFailed
	default:
		return booking.BookingStatusPending
	}
}

func (s *TransactionServiceImpl) GetTransactionByOrderCode(ctx context.Context, orderCode int) (*model.Transaction, error) {
	transaction, err := s.transactionRepo.GetByOrderCode(ctx, orderCode)
	if err != nil {
		return nil, ginext.NewNotFoundError("transaction not found")
	}
	return transaction, nil
}

func (s *TransactionServiceImpl) GetTransactionByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error) {
	transaction, err := s.transactionRepo.GetByBookingID(ctx, bookingID)
	if err != nil {
		return nil, ginext.NewNotFoundError("transaction not found")
	}
	return transaction, nil
}

func (s *TransactionServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ginext.NewNotFoundError("transaction not found")
	}
	return transaction, nil
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
