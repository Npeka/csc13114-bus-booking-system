package service

import (
	"context"
	"testing"
	"time"

	client_mocks "bus-booking/payment-service/internal/client/mocks"
	"bus-booking/payment-service/internal/model"
	repo_mocks "bus-booking/payment-service/internal/repository/mocks"
	service_mocks "bus-booking/payment-service/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/payOSHQ/payos-lib-golang/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	assert.NotNil(t, service)
	assert.IsType(t, &TransactionServiceImpl{}, service)
}

func TestGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	transactionID := uuid.New()

	transaction := &model.Transaction{
		BaseModel:     model.BaseModel{ID: transactionID},
		BookingID:     uuid.New(),
		UserID:        uuid.New(),
		Amount:        100000,
		Status:        model.TransactionStatusPaid,
		PaymentMethod: model.PaymentMethodPayOS,
	}

	mockTransactionRepo.EXPECT().
		GetByID(ctx, transactionID).
		Return(transaction, nil).
		Times(1)

	result, err := service.GetByID(ctx, transactionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, model.TransactionStatusPaid, result.Status)
	assert.Equal(t, 100000, result.Amount)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	transactionID := uuid.New()

	mockTransactionRepo.EXPECT().
		GetByID(ctx, transactionID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetByID(ctx, transactionID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetByBookingID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	bookingID := uuid.New()

	transaction := &model.Transaction{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		BookingID:     bookingID,
		UserID:        uuid.New(),
		Amount:        200000,
		Status:        model.TransactionStatusPending,
		PaymentMethod: model.PaymentMethodPayOS,
	}

	mockTransactionRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(transaction, nil).
		Times(1)

	result, err := service.GetByBookingID(ctx, bookingID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.BookingID)
	assert.Equal(t, 200000, result.Amount)
}

func TestGetByBookingID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	bookingID := uuid.New()

	mockTransactionRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetByBookingID(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()

	status := model.TransactionStatusPaid
	query := &model.TransactionListQuery{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
		Status: &status,
	}

	transactions := []*model.Transaction{
		{
			BaseModel:     model.BaseModel{ID: uuid.New()},
			BookingID:     uuid.New(),
			UserID:        uuid.New(),
			Amount:        100000,
			Status:        model.TransactionStatusPaid,
			PaymentMethod: model.PaymentMethodPayOS,
		},
		{
			BaseModel:     model.BaseModel{ID: uuid.New()},
			BookingID:     uuid.New(),
			UserID:        uuid.New(),
			Amount:        150000,
			Status:        model.TransactionStatusPaid,
			PaymentMethod: model.PaymentMethodPayOS,
		},
	}

	mockTransactionRepo.EXPECT().
		GetList(ctx, query).
		Return(transactions, int64(2), nil).
		Times(1)

	result, total, err := service.GetList(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, model.TransactionStatusPaid, result[0].Status)
	assert.Equal(t, 100000, result[0].Amount)
}

func TestGetList_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()

	query := &model.TransactionListQuery{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	mockTransactionRepo.EXPECT().
		GetList(ctx, query).
		Return([]*model.Transaction{}, int64(0), nil).
		Times(1)

	result, total, err := service.GetList(ctx, query)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
}

func TestGetStats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()

	expectedStats := &model.TransactionStats{
		TotalTransactions:  100,
		TotalIn:            8000000,
		TotalOut:           1500000,
		PendingRefunds:     1000000,
		PendingRefundCount: 15,
	}

	mockTransactionRepo.EXPECT().
		GetStats(ctx).
		Return(expectedStats, nil).
		Times(1)

	result, err := service.GetStats(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 100, result.TotalTransactions)
	assert.Equal(t, 8000000, result.TotalIn)
	assert.Equal(t, 1500000, result.TotalOut)
	assert.Equal(t, 1000000, result.PendingRefunds)
	assert.Equal(t, 15, result.PendingRefundCount)
}

func TestGetStats_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()

	mockTransactionRepo.EXPECT().
		GetStats(ctx).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetStats(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCancel_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	transactionID := uuid.New()
	paymentLinkID := "payos-payment-link-123"

	transaction := &model.Transaction{
		BaseModel:     model.BaseModel{ID: transactionID},
		BookingID:     uuid.New(),
		UserID:        uuid.New(),
		Amount:        100000,
		Status:        model.TransactionStatusPending,
		PaymentMethod: model.PaymentMethodPayOS,
		PaymentLinkID: paymentLinkID,
	}

	cancellationReason := "Booking cancelled by user"
	paymentLink := &payos.PaymentLink{
		Status: payos.PaymentLinkStatusCancelled,
	}

	mockTransactionRepo.EXPECT().
		GetByID(ctx, transactionID).
		Return(transaction, nil).
		Times(1)

	mockPayOSService.EXPECT().
		CancelPaymentLink(ctx, paymentLinkID, &cancellationReason).
		Return(paymentLink, nil).
		Times(1)

	mockPayOSService.EXPECT().
		ToTransactionStatus(payos.PaymentLinkStatusCancelled).
		Return(model.TransactionStatusCancelled).
		Times(1)

	mockTransactionRepo.EXPECT().
		UpdateTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, tx *model.Transaction) error {
			assert.Equal(t, model.TransactionStatusCancelled, tx.Status)
			return nil
		}).
		Times(1)

	result, err := service.Cancel(ctx, transactionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, model.TransactionStatusCancelled, result.Status)
}

func TestCancel_AlreadyPaid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	transactionID := uuid.New()

	transaction := &model.Transaction{
		BaseModel:     model.BaseModel{ID: transactionID},
		BookingID:     uuid.New(),
		UserID:        uuid.New(),
		Amount:        100000,
		Status:        model.TransactionStatusPaid, // Already paid
		PaymentMethod: model.PaymentMethodPayOS,
	}

	mockTransactionRepo.EXPECT().
		GetByID(ctx, transactionID).
		Return(transaction, nil).
		Times(1)

	result, err := service.Cancel(ctx, transactionID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot cancel")
}

func TestCancel_TransactionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	transactionID := uuid.New()

	mockTransactionRepo.EXPECT().
		GetByID(ctx, transactionID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.Cancel(ctx, transactionID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()
	transactionID := uuid.New()

	req := &model.CreateTransactionRequest{
		ID:            transactionID,
		BookingID:     bookingID,
		Amount:        100000,
		Currency:      model.CurrencyVND,
		PaymentMethod: model.PaymentMethodPayOS,
		Description:   "Test payment",
		ExpiresAt:     time.Now().Add(15 * time.Minute),
	}

	payosResponse := &payos.CreatePaymentLinkResponse{
		OrderCode:     123456,
		PaymentLinkId: "payos-payment-link-123",
		CheckoutUrl:   "https://checkout.url",
		QrCode:        "qr-code-data",
		Status:        payos.PaymentLinkStatusPending,
	}

	mockPayOSService.EXPECT().
		CreatePaymentLink(ctx, gomock.Any()).
		Return(payosResponse, nil).
		Times(1)

	mockPayOSService.EXPECT().
		ToTransactionStatus(payos.PaymentLinkStatusPending).
		Return(model.TransactionStatusPending).
		Times(1)

	mockTransactionRepo.EXPECT().
		CreateTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, tx *model.Transaction) error {
			assert.Equal(t, transactionID, tx.ID)
			assert.Equal(t, bookingID, tx.BookingID)
			assert.Equal(t, userID, tx.UserID)
			assert.Equal(t, 100000, tx.Amount)
			assert.Equal(t, model.TransactionStatusPending, tx.Status)
			assert.Equal(t, int64(123456), tx.OrderCode)
			assert.Equal(t, "payos-payment-link-123", tx.PaymentLinkID)
			return nil
		}).
		Times(1)

	result, err := service.Create(ctx, req, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.BookingID)
	assert.Equal(t, "https://checkout.url", result.CheckoutURL)
}

func TestCreate_PayOSError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBookingClient := client_mocks.NewMockBookingClient(ctrl)
	mockPayOSService := service_mocks.NewMockPayOSService(ctrl)

	service := NewTransactionService(mockTransactionRepo, mockBookingClient, mockPayOSService)

	ctx := context.Background()
	userID := uuid.New()

	req := &model.CreateTransactionRequest{
		ID:          uuid.New(),
		BookingID:   uuid.New(),
		Amount:      100000,
		Description: "Test payment",
		ExpiresAt:   time.Now().Add(15 * time.Minute),
	}

	mockPayOSService.EXPECT().
		CreatePaymentLink(ctx, gomock.Any()).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.Create(ctx, req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create payment link")
}
