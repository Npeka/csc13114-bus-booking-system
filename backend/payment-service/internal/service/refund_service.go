package service

import (
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/repository"
	"bus-booking/shared/ginext"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RefundService interface {
	CreateRefund(ctx context.Context, req *model.RefundRequest, userID uuid.UUID) (*model.RefundResponse, error)
	GetRefundByBookingID(ctx context.Context, bookingID uuid.UUID, userID uuid.UUID) (*model.RefundResponse, error)
	ListRefunds(ctx context.Context, query *model.RefundListQuery) ([]*model.RefundResponse, int64, error)
	UpdateRefundStatus(ctx context.Context, transactionID uuid.UUID, status model.RefundStatus, adminID uuid.UUID) error
	ExportRefundsToExcel(ctx context.Context, refundIDs []uuid.UUID) ([]byte, error)
}

type RefundServiceImpl struct {
	refundRepo       repository.RefundRepository // NEW - dedicated refund repository
	transactionRepo  repository.TransactionRepository
	bankAccountRepo  repository.BankAccountRepository
	constantsService ConstantsService
	excelService     ExcelService
}

func NewRefundService(
	refundRepo repository.RefundRepository, // NEW first param
	transactionRepo repository.TransactionRepository,
	bankAccountRepo repository.BankAccountRepository,
	constantsService ConstantsService,
	excelService ExcelService,
) RefundService {
	return &RefundServiceImpl{
		refundRepo:       refundRepo,
		transactionRepo:  transactionRepo,
		bankAccountRepo:  bankAccountRepo,
		constantsService: constantsService,
		excelService:     excelService,
	}
}

func (s *RefundServiceImpl) CreateRefund(ctx context.Context, req *model.RefundRequest, userID uuid.UUID) (*model.RefundResponse, error) {
	// Get original payment transaction
	originalTx, err := s.transactionRepo.GetByBookingID(ctx, req.BookingID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get original transaction")
		return nil, ginext.NewNotFoundError("original transaction not found")
	}

	// Verify it's a paid transaction
	if originalTx.Status != model.TransactionStatusPaid {
		return nil, ginext.NewBadRequestError("cannot refund unpaid transaction")
	}

	// Verify user owns this transaction
	if originalTx.UserID != userID {
		return nil, ginext.NewForbiddenError("you don't own this transaction")
	}

	// Check if refund already exists for this booking
	existingRefund, err := s.refundRepo.GetByBookingID(ctx, req.BookingID)
	if err == nil && existingRefund != nil {
		return nil, ginext.NewConflictError("refund already exists for this booking")
	}

	// Check if user has a primary bank account
	_, err = s.bankAccountRepo.GetPrimaryBankAccount(ctx, userID)
	if err != nil {
		return nil, ginext.NewBadRequestError("you must add a bank account before requesting refund")
	}

	// Validate refund amount
	if req.RefundAmount > originalTx.Amount {
		return nil, ginext.NewBadRequestError("refund amount cannot exceed original amount")
	}

	// Create refund entity
	refund := &model.Refund{
		BookingID:     req.BookingID,
		TransactionID: originalTx.ID,
		UserID:        userID,
		RefundAmount:  req.RefundAmount,
		RefundStatus:  model.RefundStatusPending,
		RefundReason:  req.Reason,
	}

	if err := s.refundRepo.Create(ctx, refund); err != nil {
		log.Error().Err(err).Msg("Failed to create refund")
		return nil, ginext.NewInternalServerError("failed to create refund")
	}

	return &model.RefundResponse{
		ID:                    refund.ID,
		CreatedAt:             refund.CreatedAt,
		UpdatedAt:             refund.UpdatedAt,
		BookingID:             refund.BookingID,
		UserID:                refund.UserID,
		RefundAmount:          refund.RefundAmount,
		RefundStatus:          refund.RefundStatus,
		RefundReason:          refund.RefundReason,
		OriginalTransactionID: refund.TransactionID, // For compatibility
	}, nil
}

func (s *RefundServiceImpl) GetRefundByBookingID(ctx context.Context, bookingID uuid.UUID, userID uuid.UUID) (*model.RefundResponse, error) {
	refund, err := s.refundRepo.GetByBookingID(ctx, bookingID)
	if err != nil {
		log.Error().Err(err).Msg("Refund not found")
		return nil, ginext.NewNotFoundError("refund not found for this booking")
	}

	// Verify user owns this refund
	if refund.UserID != userID {
		return nil, ginext.NewForbiddenError("you don't own this refund")
	}

	return s.toRefundResponse(refund), nil
}

func (s *RefundServiceImpl) ListRefunds(ctx context.Context, query *model.RefundListQuery) ([]*model.RefundResponse, int64, error) {
	refunds, total, err := s.refundRepo.List(ctx, query)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("failed to list refunds")
	}

	responses := make([]*model.RefundResponse, len(refunds))
	for i, refund := range refunds {
		response := s.toRefundResponse(refund)

		// Populate bank account info for admin view
		bankAccount, err := s.bankAccountRepo.GetPrimaryBankAccount(ctx, refund.UserID)
		if err == nil && bankAccount != nil {
			response.BankCode = bankAccount.BankCode
			response.BankName = s.getBankName(ctx, bankAccount.BankCode)
			response.AccountNumber = bankAccount.AccountNumber
			response.AccountHolder = bankAccount.AccountHolder
		}

		responses[i] = response
	}

	return responses, total, nil
}

func (s *RefundServiceImpl) UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, status model.RefundStatus, adminID uuid.UUID) error {
	// Validate status transition
	if status != model.RefundStatusProcessing &&
		status != model.RefundStatusCompleted &&
		status != model.RefundStatusRejected {
		return ginext.NewBadRequestError("invalid refund status")
	}

	// Get refund
	refund, err := s.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		log.Error().Err(err).Msg("Refund not found")
		return ginext.NewNotFoundError("refund not found")
	}

	// Update fields
	refund.RefundStatus = status
	refund.ProcessedBy = &adminID
	now := time.Now()
	refund.ProcessedAt = &now

	if err := s.refundRepo.Update(ctx, refund); err != nil {
		log.Error().Err(err).Msg("Failed to update refund status")
		return ginext.NewInternalServerError("failed to update refund status")
	}

	return nil
}

func (s *RefundServiceImpl) ExportRefundsToExcel(ctx context.Context, refundIDs []uuid.UUID) ([]byte, error) {
	// Get refunds by IDs
	refunds, err := s.refundRepo.ListByIDs(ctx, refundIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get refunds")
		return nil, ginext.NewInternalServerError("failed to get refunds")
	}

	if len(refunds) == 0 {
		return nil, ginext.NewNotFoundError("no refunds found")
	}

	// Get user bank accounts for each refund
	exportItems := make([]*model.RefundExportItem, 0, len(refunds))
	for _, refund := range refunds {
		// Get user's primary bank account
		bankAccount, err := s.bankAccountRepo.GetPrimaryBankAccount(ctx, refund.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", refund.UserID.String()).Msg("User has no bank account")
			// Skip this refund if no bank account
			continue
		}

		// Get bank name
		bankName := s.getBankName(ctx, bankAccount.BankCode)

		// Generate booking reference from BookingID
		bookingRef := refund.BookingID.String()[:8]

		exportItems = append(exportItems, &model.RefundExportItem{
			BookingReference: bookingRef,
			UserName:         bankAccount.AccountHolder,
			BankCode:         bankAccount.BankCode,
			BankName:         bankName,
			AccountNumber:    bankAccount.AccountNumber,
			AccountHolder:    bankAccount.AccountHolder,
			RefundAmount:     refund.RefundAmount,
			Reason:           refund.RefundReason,
			CreatedDate:      refund.CreatedAt,
		})
	}

	if len(exportItems) == 0 {
		return nil, ginext.NewBadRequestError("no refunds with valid bank accounts")
	}

	// Generate Excel
	excelData, err := s.excelService.GenerateRefundExcel(exportItems)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate Excel")
		return nil, ginext.NewInternalServerError("failed to generate Excel file")
	}

	return excelData, nil
}

// Helper functions
func (s *RefundServiceImpl) toRefundResponse(refund *model.Refund) *model.RefundResponse {
	return &model.RefundResponse{
		ID:                    refund.ID,
		CreatedAt:             refund.CreatedAt,
		UpdatedAt:             refund.UpdatedAt,
		BookingID:             refund.BookingID,
		UserID:                refund.UserID,
		RefundAmount:          refund.RefundAmount,
		RefundStatus:          refund.RefundStatus,
		RefundReason:          refund.RefundReason,
		OriginalTransactionID: refund.TransactionID,
		ProcessedBy:           refund.ProcessedBy,
		ProcessedAt:           refund.ProcessedAt,
	}
}

func (s *RefundServiceImpl) getBankName(ctx context.Context, bankCode string) string {
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
