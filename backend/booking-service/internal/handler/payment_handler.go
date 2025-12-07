package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/rs/zerolog/log"
)

type PaymentHandler interface {
	GetPaymentMethods(r *ginext.Request) (*ginext.Response, error)
	ProcessPayment(r *ginext.Request) (*ginext.Response, error)
}

type PaymentHandlerImpl struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) PaymentHandler {
	return &PaymentHandlerImpl{
		paymentService: paymentService,
	}
}

// GetPaymentMethods godoc
// @Summary Get payment methods
// @Description Get all available payment methods
// @Tags payment
// @Produce json
// @Success 200 {object} ginext.Response{data=[]model.PaymentMethodResponse}
// @Failure 500 {object} ginext.Response
// @Router /api/v1/payment/methods [get]
func (h *PaymentHandlerImpl) GetPaymentMethods(r *ginext.Request) (*ginext.Response, error) {
	methods, err := h.paymentService.GetPaymentMethods(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get payment methods")
		return nil, err
	}

	return ginext.NewSuccessResponse(methods), nil
}

// ProcessPayment godoc
// @Summary Process payment
// @Description Process payment for a booking
// @Tags payment
// @Accept json
// @Produce json
// @Param request body model.ProcessPaymentRequest true "Payment processing request"
// @Success 200 {object} ginext.Response{data=model.PaymentResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/payment/process [post]
func (h *PaymentHandlerImpl) ProcessPayment(r *ginext.Request) (*ginext.Response, error) {
	var req model.ProcessPaymentRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	payment, err := h.paymentService.ProcessPayment(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("failed to process payment")
		return nil, err
	}

	return ginext.NewSuccessResponse(payment), nil
}
