package handler

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/service"
	sharedcontext "bus-booking/shared/context"
	"bus-booking/shared/ginext"
)

type TransactionHandler interface {
	CreatePaymentLink(r *ginext.Request) (*ginext.Response, error)
	HandlePaymentWebhook(r *ginext.Request) (*ginext.Response, error)
	CancelPayment(r *ginext.Request) (*ginext.Response, error)
	GetByID(r *ginext.Request) (*ginext.Response, error)
}

type TransactionHandlerImpl struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) TransactionHandler {
	return &TransactionHandlerImpl{
		service: service,
	}
}

// CreatePaymentLink godoc
// @Summary Create a payment link
// @Description Create a payment link via PayOS for a booking
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body model.CreatePaymentLinkRequest true "Payment creation request"
// @Success 201 {object} ginext.Response{data=model.TransactionResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/payment-link [post]
func (h *TransactionHandlerImpl) CreatePaymentLink(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	var req model.CreatePaymentLinkRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	resp, err := h.service.CreatePaymentLink(r.GinCtx.Request.Context(), &req, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment link")
		return nil, ginext.NewInternalServerError("Failed to create payment link")
	}

	return ginext.NewSuccessResponse(resp), nil
}

// HandlePaymentWebhook godoc
// @Summary Handle PayOS webhook
// @Description Handle payment webhook notification from PayOS
// @Tags transactions
// @Accept json
// @Produce json
// @Param webhook body model.PaymentWebhookData true "Webhook payload"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/webhook [post]
func (h *TransactionHandlerImpl) HandlePaymentWebhook(r *ginext.Request) (*ginext.Response, error) {
	log.Info().Msg("Webhook handler started")

	var webhookData map[string]interface{}
	if err := r.GinCtx.ShouldBindJSON(webhookData); err != nil {
		log.Error().Err(err).Msg("JSON binding failed - invalid JSON format")
		return nil, ginext.NewBadRequestError("Invalid webhook data")
	}

	log.Info().Interface("webhookData", webhookData).Msg("Successfully parsed webhook data")

	err := h.service.HandlePaymentWebhook(r.GinCtx.Request.Context(), webhookData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process webhook in service layer")
		return nil, err
	}

	log.Info().Msg("Webhook processed successfully")
	return ginext.NewSuccessResponse("Webhook processed successfully"), nil
}

// GetByID godoc
// @Summary Get transaction by ID
// @Description Retrieve transaction details by transaction ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} ginext.Response{data=model.TransactionResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandlerImpl) GetByID(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid transaction ID")
	}

	transaction, err := h.service.GetByID(r.GinCtx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Transaction not found")
		return nil, ginext.NewNotFoundError("Transaction not found")
	}

	return ginext.NewSuccessResponse(transaction), nil
}

// CancelPayment godoc
// @Summary Cancel a payment
// @Description Cancel a payment transaction and PayOS payment link
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/{id}/cancel [post]
func (h *TransactionHandlerImpl) CancelPayment(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid transaction ID")
	}

	transaction, err := h.service.CancelPayment(r.GinCtx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Failed to cancel payment")
		return nil, err
	}

	return ginext.NewSuccessResponse(transaction), nil
}
