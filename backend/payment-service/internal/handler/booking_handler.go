package handler

import (
	"github.com/rs/zerolog/log"

	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/service"
	"bus-booking/shared/ginext"
)

type TransactionHandler interface {
	CreateTransaction(r *ginext.Request) (*ginext.Response, error)
}

type TransactionHandlerImpl struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) TransactionHandler {
	return &TransactionHandlerImpl{
		service: service,
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction for a booking
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body model.CreateTransactionRequest true "Transaction creation request"
// @Success 201 {object} model.TransactionResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /transactions [post]
func (h *TransactionHandlerImpl) CreateTransaction(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateTransactionRequest
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.service.CreateTransaction(r.GinCtx.Request.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Failed to create transaction")
		return nil, ginext.NewInternalServerError("Failed to create transaction")
	}

	return ginext.NewSuccessResponse(nil, "Transaction created successfully"), nil
}
