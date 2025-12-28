package service

import (
	"context"
	"testing"

	client_mocks "bus-booking/payment-service/internal/client/mocks"
	"bus-booking/payment-service/internal/model"
	repo_mocks "bus-booking/payment-service/internal/repository/mocks"
	service_mocks "bus-booking/payment-service/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/payOSHQ/payos-lib-golang/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebhook_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	bookingID := uuid.New()
	paymentLinkID := "payos-payment-link-123"
	orderCode := 123456

	webhookMap := map[string]interface{}{
		"code": "00",
		"desc": "success",
		"data": map[string]interface{}{
			"orderCode":     orderCode,
			"paymentLinkId": paymentLinkID,
		},
	}

	webhookData := model.PaymentWebhookData{
		Code:    "00",
		Desc:    "success",
		Success: true,
		Data: model.PaymentWebhookDetails{
			OrderCode:           orderCode,
			PaymentLinkID:       paymentLinkID,
			Reference:           "REF123",
			TransactionDateTime: "2024-01-15 10:30:00",
		},
	}

	transaction := &model.Transaction{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		BookingID:     bookingID,
		UserID:        uuid.New(),
		Amount:        100000,
		Status:        model.TransactionStatusPending,
		PaymentMethod: model.PaymentMethodPayOS,
		PaymentLinkID: paymentLinkID,
		OrderCode:     int64(orderCode),
	}

	paymentLink := &payos.PaymentLink{
		Status: payos.PaymentLinkStatusPaid,
	}

	// Mock webhook verification
	mockPayOSService.EXPECT().
		VerifyWebhook(ctx, webhookMap).
		Return(nil).
		Times(1)

	// Mock payment link retrieval
	mockPayOSService.EXPECT().
		GetPaymentLink(gomock.Any(), paymentLinkID).
		Return(paymentLink, nil).
		Times(1)

	// Mock transaction retrieval by webhook data
	mockTransactionRepo.EXPECT().
		GetByWebhookData(gomock.Any(), orderCode, paymentLinkID).
		Return(transaction, nil).
		Times(1)

	// Mock status conversion
	mockPayOSService.EXPECT().
		ToTransactionStatus(payos.PaymentLinkStatusPaid).
		Return(model.TransactionStatusPaid).
		Times(1)

	// Mock transaction update
	mockTransactionRepo.EXPECT().
		UpdateTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, tx *model.Transaction) error {
			assert.Equal(t, model.TransactionStatusPaid, tx.Status)
			assert.Equal(t, "REF123", tx.Reference)
			assert.NotNil(t, tx.TransactionTime)
			return nil
		}).
		Times(1)

	// Mock booking status update
	mockBookingClient.EXPECT().
		UpdateBookingStatus(ctx, gomock.Any(), bookingID).
		Return(nil).
		Times(1)

	err := service.HandleWebhook(ctx, webhookMap, webhookData)

	assert.NoError(t, err)
}

func TestHandleWebhook_VerificationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()

	webhookMap := map[string]interface{}{
		"code": "00",
	}

	webhookData := model.PaymentWebhookData{
		Code: "00",
	}

	mockPayOSService.EXPECT().
		VerifyWebhook(ctx, webhookMap).
		Return(assert.AnError).
		Times(1)

	err := service.HandleWebhook(ctx, webhookMap, webhookData)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid webhook signature")
}

func TestHandleWebhook_TransactionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	paymentLinkID := "payos-payment-link-123"
	orderCode := 123456

	webhookMap := map[string]interface{}{
		"code": "00",
	}

	webhookData := model.PaymentWebhookData{
		Code: "00",
		Data: model.PaymentWebhookDetails{
			OrderCode:     orderCode,
			PaymentLinkID: paymentLinkID,
		},
	}

	paymentLink := &payos.PaymentLink{
		Status: payos.PaymentLinkStatusPaid,
	}

	mockPayOSService.EXPECT().
		VerifyWebhook(ctx, webhookMap).
		Return(nil).
		Times(1)

	mockPayOSService.EXPECT().
		GetPaymentLink(gomock.Any(), paymentLinkID).
		Return(paymentLink, nil).
		Times(1)

	mockTransactionRepo.EXPECT().
		GetByWebhookData(gomock.Any(), orderCode, paymentLinkID).
		Return(nil, assert.AnError).
		Times(1)

	err := service.HandleWebhook(ctx, webhookMap, webhookData)

	assert.Error(t, err)
}
