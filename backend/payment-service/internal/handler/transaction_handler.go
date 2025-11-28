package handler

import (
	"strconv"

	"github.com/rs/zerolog/log"

	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/service"
	"bus-booking/shared/ginext"
)

type TransactionHandler interface {
	CreateTransaction(r *ginext.Request) (*ginext.Response, error)
	CreatePaymentLink(r *ginext.Request) (*ginext.Response, error)
	HandlePaymentWebhook(r *ginext.Request) (*ginext.Response, error)
	HandlePaymentReturn(r *ginext.Request) (*ginext.Response, error)
	HandlePaymentCancel(r *ginext.Request) (*ginext.Response, error)
	GetTransactionByOrderCode(r *ginext.Request) (*ginext.Response, error)
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
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transactions [post]
func (h *TransactionHandlerImpl) CreateTransaction(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateTransactionRequest
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	resp, err := h.service.CreateTransaction(r.GinCtx.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create transaction")
		return nil, ginext.NewInternalServerError("Failed to create transaction")
	}

	return ginext.NewSuccessResponse(resp, "Transaction created successfully"), nil
}

// CreatePaymentLink godoc
// @Summary Create a payment link
// @Description Create a payment link via PayOS for a booking
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body model.CreateTransactionRequest true "Payment creation request"
// @Success 201 {object} ginext.Response{data=model.TransactionResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/payment-link [post]
func (h *TransactionHandlerImpl) CreatePaymentLink(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateTransactionRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	resp, err := h.service.CreatePaymentLink(r.GinCtx.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment link")
		return nil, ginext.NewInternalServerError("Failed to create payment link")
	}

	return ginext.NewSuccessResponse(resp, "Payment link created successfully"), nil
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
	var webhookData model.PaymentWebhookData
	if err := r.GinCtx.ShouldBindJSON(&webhookData); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid webhook data")
	}

	err := h.service.HandlePaymentWebhook(r.GinCtx.Request.Context(), &webhookData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process webhook")
		return nil, ginext.NewInternalServerError("Failed to process webhook")
	}

	return ginext.NewSuccessResponse(nil, "Webhook processed successfully"), nil
}

// HandlePaymentReturn godoc
// @Summary Handle payment return
// @Description Handle return from PayOS payment page (success)
// @Tags transactions
// @Accept json
// @Produce json
// @Param orderCode query int true "Order code"
// @Param status query string false "Payment status"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/return [get]
func (h *TransactionHandlerImpl) HandlePaymentReturn(r *ginext.Request) (*ginext.Response, error) {
	orderCodeStr := r.GinCtx.Query("orderCode")
	if orderCodeStr == "" {
		return nil, ginext.NewBadRequestError("Order code is required")
	}

	orderCode, err := strconv.ParseInt(orderCodeStr, 10, 64)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid order code")
	}

	// Confirm payment with PayOS
	err = h.service.ConfirmPayment(r.GinCtx.Request.Context(), orderCode)
	if err != nil {
		log.Error().Err(err).Int64("order_code", orderCode).Msg("Failed to confirm payment")
		return nil, ginext.NewInternalServerError("Failed to confirm payment")
	}

	// Get updated transaction
	transaction, err := h.service.GetTransactionByOrderCode(r.GinCtx.Request.Context(), orderCode)
	if err != nil {
		return nil, ginext.NewInternalServerError("Failed to get transaction")
	}

	return ginext.NewSuccessResponse(map[string]interface{}{
		"order_code": orderCode,
		"status":     transaction.Status,
		"booking_id": transaction.BookingID,
	}, "Payment confirmed successfully"), nil
}

// HandlePaymentCancel godoc
// @Summary Handle payment cancellation
// @Description Handle cancel from PayOS payment page
// @Tags transactions
// @Accept json
// @Produce json
// @Param orderCode query int true "Order code"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/cancel [get]
func (h *TransactionHandlerImpl) HandlePaymentCancel(r *ginext.Request) (*ginext.Response, error) {
	orderCodeStr := r.GinCtx.Query("orderCode")
	if orderCodeStr == "" {
		return nil, ginext.NewBadRequestError("Order code is required")
	}

	orderCode, err := strconv.ParseInt(orderCodeStr, 10, 64)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid order code")
	}

	// Cancel payment
	err = h.service.CancelPayment(r.GinCtx.Request.Context(), orderCode, "User cancelled payment")
	if err != nil {
		log.Error().Err(err).Int64("order_code", orderCode).Msg("Failed to cancel payment")
		// Don't return error - just log it
	}

	return ginext.NewSuccessResponse(map[string]interface{}{
		"order_code": orderCode,
		"status":     "CANCELLED",
	}, "Payment cancelled"), nil
}

// GetTransactionByOrderCode godoc
// @Summary Get transaction by order code
// @Description Retrieve transaction details by PayOS order code
// @Tags transactions
// @Produce json
// @Param order_code path int true "Order code"
// @Success 200 {object} ginext.Response{data=model.TransactionResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/transactions/{order_code} [get]
func (h *TransactionHandlerImpl) GetTransactionByOrderCode(r *ginext.Request) (*ginext.Response, error) {
	orderCodeStr := r.GinCtx.Param("order_code")
	orderCode, err := strconv.ParseInt(orderCodeStr, 10, 64)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid order code")
	}

	transaction, err := h.service.GetTransactionByOrderCode(r.GinCtx.Request.Context(), orderCode)
	if err != nil {
		log.Error().Err(err).Int64("order_code", orderCode).Msg("Transaction not found")
		return nil, ginext.NewNotFoundError("Transaction not found")
	}

	return ginext.NewSuccessResponse(transaction, "Transaction retrieved successfully"), nil
}
