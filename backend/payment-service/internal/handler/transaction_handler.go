package handler

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/service"
	sharedcontext "bus-booking/shared/context"
	"bus-booking/shared/ginext"
)

type TransactionHandler interface {
	GetList(r *ginext.Request) (*ginext.Response, error)
	GetByID(r *ginext.Request) (*ginext.Response, error)
	GetStats(r *ginext.Request) (*ginext.Response, error)

	Create(r *ginext.Request) (*ginext.Response, error)
	Cancel(r *ginext.Request) (*ginext.Response, error)
	HandleWebhook(r *ginext.Request) (*ginext.Response, error)
}

type TransactionHandlerImpl struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) TransactionHandler {
	return &TransactionHandlerImpl{
		service: service,
	}
}

// GetList godoc
// @Summary List all transactions (Admin)
// @Description List all transactions with filters
// @Tags admin
// @Accept json
// @Produce json
// @Param transaction_type query string false "Transaction type" Enums(IN, OUT)
// @Param status query string false "Transaction status"
// @Param refund_status query string false "Refund status"
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions [get]
func (h *TransactionHandlerImpl) GetList(r *ginext.Request) (*ginext.Response, error) {
	var query model.TransactionListQuery
	if err := r.GinCtx.ShouldBindQuery(&query); err != nil {
		log.Debug().Err(err).Msg("Query binding failed")
		return nil, ginext.NewBadRequestError("Invalid query parameters")
	}

	// Normalize defaults
	query.Normalize()

	transactions, total, err := h.service.GetList(r.Context(), &query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list transactions")
		return nil, err
	}

	return ginext.NewPaginatedResponse(transactions, query.Page, query.PageSize, total), nil
}

// GetStats godoc
// @Summary Get transaction stats (Admin)
// @Description Get transaction stats with filters
// @Tags admin
// @Accept json
// @Produce json
// @Param transaction_type query string false "Transaction type" Enums(IN, OUT)
// @Param status query string false "Transaction status"
// @Param refund_status query string false "Refund status"
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/stats [get]
func (h *TransactionHandlerImpl) GetStats(r *ginext.Request) (*ginext.Response, error) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transaction stats")
		return nil, err
	}

	return ginext.NewSuccessResponse(stats), nil
}

// Create godoc
// @Summary Create a payment link
// @Description Create a payment link via PayOS for a booking
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body model.CreateTransactionRequest true "Payment creation request"
// @Success 201 {object} ginext.Response{data=model.TransactionResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions [post]
func (h *TransactionHandlerImpl) Create(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	var req model.CreateTransactionRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	resp, err := h.service.Create(r.GinCtx.Request.Context(), &req, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment link")
		return nil, ginext.NewInternalServerError("Failed to create payment link")
	}

	return ginext.NewSuccessResponse(resp), nil
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

// Cancel godoc
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
func (h *TransactionHandlerImpl) Cancel(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid transaction ID")
	}

	transaction, err := h.service.Cancel(r.GinCtx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Failed to cancel payment")
		return nil, err
	}

	return ginext.NewSuccessResponse(transaction), nil
}

// HandleWebhook godoc
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
func (h *TransactionHandlerImpl) HandleWebhook(r *ginext.Request) (*ginext.Response, error) {
	log.Info().Msg("Webhook handler started")

	// Read the raw body first (needed for signature verification)
	bodyBytes, err := io.ReadAll(r.GinCtx.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		return nil, ginext.NewBadRequestError("Invalid request body")
	}
	defer func() {
		if err := r.GinCtx.Request.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close request body")
		}
	}()

	// Parse to map for verification
	var webhookMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &webhookMap); err != nil {
		log.Error().Err(err).Msg("JSON parsing failed - invalid JSON format")
		return nil, ginext.NewBadRequestError("Invalid webhook data")
	}

	// Parse to struct for processing
	var webhookData model.PaymentWebhookData
	if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
		log.Error().Err(err).Msg("JSON parsing failed - invalid webhook structure")
		return nil, ginext.NewBadRequestError("Invalid webhook data")
	}

	log.Info().Interface("webhookData", webhookData).Msg("Successfully parsed webhook data")

	err = h.service.HandleWebhook(r.GinCtx.Request.Context(), webhookMap, webhookData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process webhook in service layer")
		return nil, err
	}

	log.Info().Msg("Webhook processed successfully")
	return ginext.NewSuccessResponse("Webhook processed successfully"), nil
}
