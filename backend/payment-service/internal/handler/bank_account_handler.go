package handler

import (
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/service"
	sharedcontext "bus-booking/shared/context"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BankAccountHandler interface {
	CreateBankAccount(r *ginext.Request) (*ginext.Response, error)
	GetMyBankAccounts(r *ginext.Request) (*ginext.Response, error)
	UpdateBankAccount(r *ginext.Request) (*ginext.Response, error)
	DeleteBankAccount(r *ginext.Request) (*ginext.Response, error)
	SetPrimaryBankAccount(r *ginext.Request) (*ginext.Response, error)
}

type BankAccountHandlerImpl struct {
	service service.BankAccountService
}

func NewBankAccountHandler(service service.BankAccountService) BankAccountHandler {
	return &BankAccountHandlerImpl{
		service: service,
	}
}

// CreateBankAccount godoc
// @Summary Create a bank account
// @Description Create a new bank account for the authenticated user
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param account body model.BankAccountRequest true "Bank account details"
// @Success 201 {object} ginext.Response{data=model.BankAccountResponse}
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bank-accounts [post]
func (h *BankAccountHandlerImpl) CreateBankAccount(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	var req model.BankAccountRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	account, err := h.service.CreateBankAccount(r.Context(), &req, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create bank account")
		return nil, err
	}

	return ginext.NewCreatedResponse(account), nil
}

// GetMyBankAccounts godoc
// @Summary Get my bank accounts
// @Description Get all bank accounts of the authenticated user
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Success 200 {object} ginext.Response{data=[]model.BankAccountResponse}
// @Failure 401 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bank-accounts [get]
func (h *BankAccountHandlerImpl) GetMyBankAccounts(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	accounts, err := h.service.GetUserBankAccounts(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get bank accounts")
		return nil, err
	}

	return ginext.NewSuccessResponse(accounts), nil
}

// UpdateBankAccount godoc
// @Summary Update a bank account
// @Description Update a bank account of the authenticated user
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Param account body model.BankAccountRequest true "Bank account details"
// @Success 200 {object} ginext.Response{data=model.BankAccountResponse}
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bank-accounts/{id} [put]
func (h *BankAccountHandlerImpl) UpdateBankAccount(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid bank account ID")
	}

	var req model.BankAccountRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	account, err := h.service.UpdateBankAccount(r.Context(), id, &req, userID)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Failed to update bank account")
		return nil, err
	}

	return ginext.NewSuccessResponse(account), nil
}

// DeleteBankAccount godoc
// @Summary Delete a bank account
// @Description Delete a bank account of the authenticated user
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bank-accounts/{id} [delete]
func (h *BankAccountHandlerImpl) DeleteBankAccount(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid bank account ID")
	}

	if err := h.service.DeleteBankAccount(r.Context(), id, userID); err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Failed to delete bank account")
		return nil, err
	}

	return ginext.NewSuccessResponse("Bank account deleted successfully"), nil
}

// SetPrimaryBankAccount godoc
// @Summary Set primary bank account
// @Description Set a bank account as the primary account for refunds
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bank-accounts/{id}/set-primary [post]
func (h *BankAccountHandlerImpl) SetPrimaryBankAccount(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid bank account ID")
	}

	if err := h.service.SetPrimaryBankAccount(r.Context(), id, userID); err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Failed to set primary bank account")
		return nil, err
	}

	return ginext.NewSuccessResponse("Primary bank account set successfully"), nil
}
