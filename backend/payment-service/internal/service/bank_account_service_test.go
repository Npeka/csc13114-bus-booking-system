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

func TestNewBankAccountService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	assert.NotNil(t, service)
	assert.IsType(t, &BankAccountServiceImpl{}, service)
}

func TestCreateBankAccount_FirstAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()

	req := &model.BankAccountRequest{
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		AccountHolder: "NGUYEN VAN A",
	}

	banks := []model.BankConstant{
		{Code: "VCB", ShortName: "Vietcombank"},
	}

	// Mock bank code validation
	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	// Mock checking existing accounts (first account)
	mockBankAccountRepo.EXPECT().
		GetUserBankAccounts(ctx, userID).
		Return([]*model.BankAccount{}, nil). // No existing accounts
		Times(1)

	// Mock account creation
	mockBankAccountRepo.EXPECT().
		CreateBankAccount(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, account *model.BankAccount) error {
			assert.Equal(t, userID, account.UserID)
			assert.Equal(t, "VCB", account.BankCode)
			assert.Equal(t, "1234567890", account.AccountNumber)
			assert.Equal(t, "NGUYEN VAN A", account.AccountHolder)
			assert.True(t, account.IsPrimary) // First account is primary
			account.ID = uuid.New()
			return nil
		}).
		Times(1)

	// Mock getting bank name for response
	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	result, err := service.CreateBankAccount(ctx, req, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1234567890", result.AccountNumber)
	assert.Equal(t, "Vietcombank", result.BankName)
	assert.True(t, result.IsPrimary)
}

func TestCreateBankAccount_SecondAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()

	req := &model.BankAccountRequest{
		BankCode:      "TCB",
		AccountNumber: "0987654321",
		AccountHolder: "NGUYEN VAN B",
	}

	banks := []model.BankConstant{
		{Code: "VCB", ShortName: "Vietcombank"},
		{Code: "TCB", ShortName: "Techcombank"},
	}

	existingAccount := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		UserID:        userID,
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		IsPrimary:     true,
	}

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	// Mock checking existing accounts (has one existing)
	mockBankAccountRepo.EXPECT().
		GetUserBankAccounts(ctx, userID).
		Return([]*model.BankAccount{existingAccount}, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		CreateBankAccount(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, account *model.BankAccount) error {
			assert.Equal(t, "TCB", account.BankCode)
			assert.False(t, account.IsPrimary) // Not primary since there's existing account
			account.ID = uuid.New()
			return nil
		}).
		Times(1)

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	result, err := service.CreateBankAccount(ctx, req, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsPrimary)
}

func TestCreateBankAccount_InvalidBankCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()

	req := &model.BankAccountRequest{
		BankCode:      "INVALID",
		AccountNumber: "1234567890",
		AccountHolder: "NGUYEN VAN A",
	}

	banks := []model.BankConstant{
		{Code: "VCB", ShortName: "Vietcombank"},
	}

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	result, err := service.CreateBankAccount(ctx, req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid bank code")
}

func TestGetUserBankAccounts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()

	accounts := []*model.BankAccount{
		{
			BaseModel:     model.BaseModel{ID: uuid.New()},
			UserID:        userID,
			BankCode:      "VCB",
			AccountNumber: "1234567890",
			IsPrimary:     true,
		},
		{
			BaseModel:     model.BaseModel{ID: uuid.New()},
			UserID:        userID,
			BankCode:      "TCB",
			AccountNumber: "0987654321",
			IsPrimary:     false,
		},
	}

	banks := []model.BankConstant{
		{Code: "VCB", ShortName: "Vietcombank"},
		{Code: "TCB", ShortName: "Techcombank"},
	}

	mockBankAccountRepo.EXPECT().
		GetUserBankAccounts(ctx, userID).
		Return(accounts, nil).
		Times(1)

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(2) // Called for each account

	result, err := service.GetUserBankAccounts(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Vietcombank", result[0].BankName)
	assert.True(t, result[0].IsPrimary)
	assert.Equal(t, "Techcombank", result[1].BankName)
	assert.False(t, result[1].IsPrimary)
}

func TestGetPrimaryBankAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()

	account := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		UserID:        userID,
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		IsPrimary:     true,
	}

	banks := []model.BankConstant{
		{Code: "VCB", ShortName: "Vietcombank"},
	}

	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID).
		Return(account, nil).
		Times(1)

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	result, err := service.GetPrimaryBankAccount(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsPrimary)
	assert.Equal(t, "Vietcombank", result.BankName)
}

func TestGetPrimaryBankAccount_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()

	mockBankAccountRepo.EXPECT().
		GetPrimaryBankAccount(ctx, userID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetPrimaryBankAccount(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no primary")
}

func TestUpdateBankAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	req := &model.BankAccountRequest{
		BankCode:      "TCB",
		AccountNumber: "9999999999",
		AccountHolder: "NGUYEN VAN C",
	}

	existingAccount := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: accountID},
		UserID:        userID,
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		AccountHolder: "NGUYEN VAN A",
		IsPrimary:     true,
	}

	banks := []model.BankConstant{
		{Code: "VCB", ShortName: "Vietcombank"},
		{Code: "TCB", ShortName: "Techcombank"},
	}

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		GetByID(ctx, accountID, userID).
		Return(existingAccount, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		UpdateBankAccount(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, account *model.BankAccount) error {
			assert.Equal(t, "TCB", account.BankCode)
			assert.Equal(t, "9999999999", account.AccountNumber)
			assert.Equal(t, "NGUYEN VAN C", account.AccountHolder)
			return nil
		}).
		Times(1)

	mockConstantsService.EXPECT().
		GetBanks(ctx).
		Return(banks, nil).
		Times(1)

	result, err := service.UpdateBankAccount(ctx, accountID, req, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Techcombank", result.BankName)
}

func TestDeleteBankAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	account := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: accountID},
		UserID:        userID,
		BankCode:      "TCB",
		AccountNumber: "0987654321",
		IsPrimary:     false, // Not primary
	}

	mockBankAccountRepo.EXPECT().
		GetByID(ctx, accountID, userID).
		Return(account, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		DeleteBankAccount(ctx, accountID, userID).
		Return(nil).
		Times(1)

	err := service.DeleteBankAccount(ctx, accountID, userID)

	assert.NoError(t, err)
}

func TestDeleteBankAccount_CannotDeletePrimaryWithOthers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	primaryAccount := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: accountID},
		UserID:        userID,
		BankCode:      "VCB",
		AccountNumber: "1234567890",
		IsPrimary:     true,
	}

	otherAccount := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: uuid.New()},
		UserID:        userID,
		BankCode:      "TCB",
		AccountNumber: "0987654321",
		IsPrimary:     false,
	}

	mockBankAccountRepo.EXPECT().
		GetByID(ctx, accountID, userID).
		Return(primaryAccount, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		GetUserBankAccounts(ctx, userID).
		Return([]*model.BankAccount{primaryAccount, otherAccount}, nil).
		Times(1)

	err := service.DeleteBankAccount(ctx, accountID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete primary")
}

func TestSetPrimaryBankAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	account := &model.BankAccount{
		BaseModel:     model.BaseModel{ID: accountID},
		UserID:        userID,
		BankCode:      "TCB",
		AccountNumber: "0987654321",
		IsPrimary:     false,
	}

	mockBankAccountRepo.EXPECT().
		GetByID(ctx, accountID, userID).
		Return(account, nil).
		Times(1)

	mockBankAccountRepo.EXPECT().
		SetPrimaryBankAccount(ctx, accountID, userID).
		Return(nil).
		Times(1)

	err := service.SetPrimaryBankAccount(ctx, accountID, userID)

	assert.NoError(t, err)
}

func TestSetPrimaryBankAccount_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankAccountRepo := repo_mocks.NewMockBankAccountRepository(ctrl)
	mockConstantsService := service_mocks.NewMockConstantsService(ctrl)

	service := NewBankAccountService(mockBankAccountRepo, mockConstantsService)

	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	mockBankAccountRepo.EXPECT().
		GetByID(ctx, accountID, userID).
		Return(nil, assert.AnError).
		Times(1)

	err := service.SetPrimaryBankAccount(ctx, accountID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
