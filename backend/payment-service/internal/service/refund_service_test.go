package service

import (
	"context"
	"testing"

	"bus-booking/payment-service/internal/model"
	repo_mocks "bus-booking/payment-service/internal/repository/mocks"
	service_mocks "bus-booking/payment-service/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRefundService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	assert.NotNil(t, service)
	assert.IsType(t, &RefundServiceImpl{}, service)
}

func TestCreateRefund_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()
	transactionID := uuid.New()

	req := &model.RefundRequest{
		BookingID:    bookingID,
		Reason:       "Không muốn đi",
		RefundAmount: 100000,
	}

	transaction := &model.Transaction{
		BaseModel: model.BaseModel{ID: transactionID},
		BookingID: bookingID,
		UserID:    userID,
		Amount:    100000,
		Status:    model.TransactionStatusPaid,
	}

	bankAccount := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		UserID:        userID,
		IsPrimary:     true,
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		AccountHolder: "NGUYEN VAN A",
	}

	// Mock transaction lookup
	mockTransactionRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(transaction, nil).
		Times(1)

	// Mock refund existence check
	mockRefundRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(nil, assert.AnError). // No existing refund
		Times(1)

	// Mock bank account check
	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID).
		Return(bankAccount, nil).
		Times(1)

	// Mock refund creation
	mockRefundRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, refund *model.Refund) error {
			assert.Equal(t, bookingID, refund.BookingID)
			assert.Equal(t, userID, refund.UserID)
			assert.Equal(t, req.RefundAmount, refund.RefundAmount)
			assert.Equal(t, req.Reason, refund.RefundReason)
			assert.Equal(t, model.RefundStatusPending, refund.RefundStatus)
			return nil
		}).
		Times(1)

	result, err := service.CreateRefund(ctx, req, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.BookingID)
	assert.Equal(t, model.RefundStatusPending, result.RefundStatus)
}

func TestCreateRefund_TransactionNotPaid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	req := &model.RefundRequest{
		BookingID:    bookingID,
		Reason:       "Test",
		RefundAmount: 100000,
	}

	transaction := &model.Transaction{
		BaseModel: model.BaseModel{ID: uuid.New()},
		BookingID: bookingID,
		UserID:    userID,
		Amount:    100000,
		Status:    model.TransactionStatusPending, // Not PAID
	}

	mockTransactionRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(transaction, nil).
		Times(1)

	result, err := service.CreateRefund(ctx, req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unpaid")
}

func TestCreateRefund_NotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	bookingID := uuid.New()

	req := &model.RefundRequest{
		BookingID:    bookingID,
		Reason:       "Test",
		RefundAmount: 100000,
	}

	transaction := &model.Transaction{
		BaseModel: model.BaseModel{ID: uuid.New()},
		BookingID: bookingID,
		UserID:    otherUserID, // Different user
		Amount:    100000,
		Status:    model.TransactionStatusPaid,
	}

	mockTransactionRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(transaction, nil).
		Times(1)

	result, err := service.CreateRefund(ctx, req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "own")
}

func TestCreateRefund_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	req := &model.RefundRequest{
		BookingID:    bookingID,
		Reason:       "Test",
		RefundAmount: 100000,
	}

	transaction := &model.Transaction{
		BaseModel: model.BaseModel{ID: uuid.New()},
		BookingID: bookingID,
		UserID:    userID,
		Amount:    100000,
		Status:    model.TransactionStatusPaid,
	}

	existingRefund := &model.Refund{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		BookingID:    bookingID,
		RefundStatus: model.RefundStatusPending,
	}

	mockTransactionRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(transaction, nil).
		Times(1)

	mockRefundRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(existingRefund, nil). // Refund exists
		Times(1)

	result, err := service.CreateRefund(ctx, req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
}

func TestCreateRefund_NoBankAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	req := &model.RefundRequest{
		BookingID:    bookingID,
		Reason:       "Test",
		RefundAmount: 100000,
	}

	transaction := &model.Transaction{
		BaseModel: model.BaseModel{ID: uuid.New()},
		BookingID: bookingID,
		UserID:    userID,
		Amount:    100000,
		Status:    model.TransactionStatusPaid,
	}

	mockTransactionRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(transaction, nil).
		Times(1)

	mockRefundRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(nil, assert.AnError).
		Times(1)

	// No bank account
	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.CreateRefund(ctx, req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "bank account")
}

func TestGetRefundByBookingID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	refund := &model.Refund{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		BookingID:    bookingID,
		UserID:       userID,
		RefundAmount: 100000,
		RefundStatus: model.RefundStatusPending,
		RefundReason: "Test",
	}

	mockRefundRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(refund, nil).
		Times(1)

	result, err := service.GetRefundByBookingID(ctx, bookingID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.BookingID)
	assert.Equal(t, model.RefundStatusPending, result.RefundStatus)
}

func TestGetRefundByBookingID_NotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	bookingID := uuid.New()

	refund := &model.Refund{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		BookingID:    bookingID,
		UserID:       otherUserID, // Different user
		RefundAmount: 100000,
		RefundStatus: model.RefundStatusPending,
	}

	mockRefundRepo.EXPECT().
		GetByBookingID(ctx, bookingID).
		Return(refund, nil).
		Times(1)

	result, err := service.GetRefundByBookingID(ctx, bookingID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "own")
}

func TestUpdateRefundStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	refundID := uuid.New()
	adminID := uuid.New()

	refund := &model.Refund{
		BaseModel:    model.BaseModel{ID: refundID},
		RefundStatus: model.RefundStatusPending,
	}

	mockRefundRepo.EXPECT().
		GetByID(ctx, refundID).
		Return(refund, nil).
		Times(1)

	mockRefundRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *model.Refund) error {
			assert.Equal(t, model.RefundStatusCompleted, r.RefundStatus)
			assert.NotNil(t, r.ProcessedBy)
			assert.Equal(t, adminID, *r.ProcessedBy)
			assert.NotNil(t, r.ProcessedAt)
			return nil
		}).
		Times(1)

	err := service.UpdateRefundStatus(ctx, refundID, model.RefundStatusCompleted, adminID)

	assert.NoError(t, err)
}

func TestUpdateRefundStatus_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	refundID := uuid.New()
	adminID := uuid.New()

	err := service.UpdateRefundStatus(ctx, refundID, model.RefundStatusPending, adminID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestListRefunds_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	userID := uuid.New()

	query := &model.RefundListQuery{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	refunds := []*model.Refund{
		{
			BaseModel:    model.BaseModel{ID: uuid.New()},
			UserID:       userID,
			RefundAmount: 100000,
			RefundStatus: model.RefundStatusPending,
		},
		{
			BaseModel:    model.BaseModel{ID: uuid.New()},
			UserID:       userID,
			RefundAmount: 200000,
			RefundStatus: model.RefundStatusCompleted,
		},
	}

	bankAccount := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		UserID:        userID,
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		AccountHolder: "NGUYEN VAN A",
	}

	banks := []model.BankConstant{
		{
			Code:      "VCB",
			ShortName: "Vietcombank",
		},
	}

	mockRefundRepo.EXPECT().
		List(ctx, query).
		Return(refunds, int64(2), nil).
		Times(1)

	// Mock bank account fetches
	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID).
		Return(bankAccount, nil).
		Times(2)

	// Mock bank name fetches
	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(2)

	result, total, err := service.ListRefunds(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, "VCB", result[0].BankCode)
	assert.Equal(t, "Vietcombank", result[0].BankName)
}

func TestExportRefundsToExcel_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)
	ctx := context.Background()
	refundID1 := uuid.New()
	refundID2 := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	bookingID1 := uuid.New()
	bookingID2 := uuid.New()

	refunds := []*model.Refund{
		{
			BaseModel:    model.BaseModel{ID: refundID1},
			UserID:       userID1,
			BookingID:    bookingID1,
			RefundAmount: 100000,
			RefundReason: "Refund reason 1",
		},
		{
			BaseModel:    model.BaseModel{ID: refundID2},
			UserID:       userID2,
			BookingID:    bookingID2,
			RefundAmount: 150000,
			RefundReason: "Refund reason 2",
		},
	}

	bankAccount1 := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		UserID:        userID1,
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		AccountHolder: "USER ONE",
		IsPrimary:     true,
	}

	bankAccount2 := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		UserID:        userID2,
		BankCode:      "TCB",
		AccountNumber: "9876543210",
		AccountHolder: "USER TWO",
		IsPrimary:     true,
	}

	banks := []model.BankConstant{
		{Code: "VCB", ShortName: "Vietcombank", Name: "Ngan hang Vietcombank"},
		{Code: "TCB", ShortName: "Techcombank", Name: "Ngan hang Techcombank"},
	}

	excelData := []byte{0x50, 0x4B, 0x03, 0x04} // ZIP signature

	mockRefundRepo.EXPECT().
		ListByIDs(ctx, []uuid.UUID{refundID1, refundID2}).
		Return(refunds, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID1).
		Return(bankAccount1, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID2).
		Return(bankAccount2, nil).
		Times(1)

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(2)

	mockExcelService.EXPECT().
		GenerateRefundExcel(gomock.Any()).
		DoAndReturn(func(items []*model.RefundExportItem) ([]byte, error) {
			assert.Len(t, items, 2)
			assert.Equal(t, "VCB", items[0].BankCode)
			assert.Equal(t, "TCB", items[1].BankCode)
			return excelData, nil
		}).
		Times(1)

	result, err := service.ExportRefundsToExcel(ctx, []uuid.UUID{refundID1, refundID2})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, excelData, result)
}

func TestExportRefundsToExcel_NoRefundsFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)

	ctx := context.Background()
	refundIDs := []uuid.UUID{uuid.New()}

	mockRefundRepo.EXPECT().
		ListByIDs(ctx, refundIDs).
		Return([]*model.Refund{}, nil).
		Times(1)

	result, err := service.ExportRefundsToExcel(ctx, refundIDs)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no refunds found")
}

func TestExportRefundsToExcel_NoBankAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := repo_mocks.NewMockRefundRepository(ctrl)
	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockTransactionRepo := repo_mocks.NewMockTransactionRepository(ctrl)
	mockExcelService := service_mocks.NewMockExcelService(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewRefundService(
		mockRefundRepo,
		mockTransactionRepo,
		mockBankAccountRepo,
		mockConstantsService,
		mockExcelService,
	)
	ctx := context.Background()
	refundID := uuid.New()
	userID := uuid.New()

	refunds := []*model.Refund{
		{
			BaseModel:    model.BaseModel{ID: refundID},
			UserID:       userID,
			BookingID:    uuid.New(),
			RefundAmount: 100000,
			RefundReason: "Refund reason",
		},
	}

	mockRefundRepo.EXPECT().
		ListByIDs(ctx, []uuid.UUID{refundID}).
		Return(refunds, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.ExportRefundsToExcel(ctx, []uuid.UUID{refundID})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no refunds with valid bank accounts")
}
