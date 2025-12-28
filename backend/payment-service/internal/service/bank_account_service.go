package service

import (
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/repository"
	"bus-booking/shared/ginext"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BankAccountService interface {
	CreateBankAccount(ctx context.Context, req *model.BankAccountRequest, userID uuid.UUID) (*model.BankAccountResponse, error)
	GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]*model.BankAccountResponse, error)
	GetPrimaryBankAccount(ctx context.Context, userID uuid.UUID) (*model.BankAccountResponse, error)
	UpdateBankAccount(ctx context.Context, accountID uuid.UUID, req *model.BankAccountRequest, userID uuid.UUID) (*model.BankAccountResponse, error)
	DeleteBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error
	SetPrimaryBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error
}

type BankAccountServiceImpl struct {
	bankAccountRepo  repository.BankAccountRepository
	constantsService ConstantsService
}

func NewBankAccountService(
	bankAccountRepo repository.BankAccountRepository,
	constantsService ConstantsService,
) BankAccountService {
	return &BankAccountServiceImpl{
		bankAccountRepo:  bankAccountRepo,
		constantsService: constantsService,
	}
}

func (s *BankAccountServiceImpl) CreateBankAccount(ctx context.Context, req *model.BankAccountRequest, userID uuid.UUID) (*model.BankAccountResponse, error) {
	// Validate bank code
	if err := s.validateBankCode(ctx, req.BankCode); err != nil {
		return nil, err
	}

	// Check if setting as first account (auto-primary)
	existingAccounts, err := s.bankAccountRepo.GetUserBankAccounts(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check existing accounts")
	}

	isPrimary := len(existingAccounts) == 0 // First account is automatically primary

	account := &model.BankAccount{
		UserID:        userID,
		BankCode:      req.BankCode,
		AccountNumber: req.AccountNumber,
		AccountHolder: req.AccountHolder,
		IsPrimary:     isPrimary,
	}

	if err := s.bankAccountRepo.CreateBankAccount(ctx, account); err != nil {
		return nil, ginext.NewInternalServerError("failed to create bank account")
	}

	return s.toResponse(ctx, account), nil
}

func (s *BankAccountServiceImpl) GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]*model.BankAccountResponse, error) {
	accounts, err := s.bankAccountRepo.GetUserBankAccounts(ctx, userID)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get bank accounts")
	}

	responses := make([]*model.BankAccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = s.toResponse(ctx, account)
	}

	return responses, nil
}

func (s *BankAccountServiceImpl) GetPrimaryBankAccount(ctx context.Context, userID uuid.UUID) (*model.BankAccountResponse, error) {
	account, err := s.bankAccountRepo.GetPrimaryBankAccount(ctx, userID)
	if err != nil {
		return nil, ginext.NewNotFoundError("no primary bank account found")
	}

	return s.toResponse(ctx, account), nil
}

func (s *BankAccountServiceImpl) UpdateBankAccount(ctx context.Context, accountID uuid.UUID, req *model.BankAccountRequest, userID uuid.UUID) (*model.BankAccountResponse, error) {
	// Validate bank code
	if err := s.validateBankCode(ctx, req.BankCode); err != nil {
		return nil, err
	}

	// Get existing account
	account, err := s.bankAccountRepo.GetByID(ctx, accountID, userID)
	if err != nil {
		return nil, ginext.NewNotFoundError("bank account not found")
	}

	// Update fields
	account.BankCode = req.BankCode
	account.AccountNumber = req.AccountNumber
	account.AccountHolder = req.AccountHolder

	if err := s.bankAccountRepo.UpdateBankAccount(ctx, account); err != nil {
		return nil, ginext.NewInternalServerError("failed to update bank account")
	}

	return s.toResponse(ctx, account), nil
}

func (s *BankAccountServiceImpl) DeleteBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	// Check if it's the primary account
	account, err := s.bankAccountRepo.GetByID(ctx, accountID, userID)
	if err != nil {
		return ginext.NewNotFoundError("bank account not found")
	}

	if account.IsPrimary {
		// Check if there are other accounts
		accounts, err := s.bankAccountRepo.GetUserBankAccounts(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get bank accounts")
		}
		if len(accounts) > 1 {
			return ginext.NewBadRequestError("cannot delete primary account. Please set another account as primary first")
		}
	}

	if err := s.bankAccountRepo.DeleteBankAccount(ctx, accountID, userID); err != nil {
		return ginext.NewInternalServerError("failed to delete bank account")
	}

	return nil
}

func (s *BankAccountServiceImpl) SetPrimaryBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	// Verify account exists and belongs to user
	if _, err := s.bankAccountRepo.GetByID(ctx, accountID, userID); err != nil {
		return ginext.NewNotFoundError("bank account not found")
	}

	if err := s.bankAccountRepo.SetPrimaryBankAccount(ctx, accountID, userID); err != nil {
		return ginext.NewInternalServerError("failed to set primary bank account")
	}

	return nil
}

// Helper functions
func (s *BankAccountServiceImpl) validateBankCode(ctx context.Context, bankCode string) error {
	banks, err := s.constantsService.GetBanks(ctx)
	if err != nil {
		return ginext.NewInternalServerError("failed to validate bank code")
	}

	for _, bank := range banks {
		if bank.Code == bankCode {
			return nil
		}
	}

	return ginext.NewBadRequestError(fmt.Sprintf("invalid bank code: %s", bankCode))
}

func (s *BankAccountServiceImpl) getBankName(ctx context.Context, bankCode string) string {
	banks, err := s.constantsService.GetBanks(ctx)
	if err != nil {
		return bankCode
	}

	for _, bank := range banks {
		if bank.Code == bankCode {
			return bank.ShortName
		}
	}

	return bankCode
}

func (s *BankAccountServiceImpl) toResponse(ctx context.Context, account *model.BankAccount) *model.BankAccountResponse {
	bankName := s.getBankName(ctx, account.BankCode)
	return account.ToResponse(bankName)
}
