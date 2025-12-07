package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type FeedbackHandler interface {
	CreateFeedback(r *ginext.Request) (*ginext.Response, error)
	GetBookingFeedback(r *ginext.Request) (*ginext.Response, error)
	GetTripFeedbacks(r *ginext.Request) (*ginext.Response, error)
}

type FeedbackHandlerImpl struct {
	feedbackService service.FeedbackService
}

func NewFeedbackHandler(feedbackService service.FeedbackService) FeedbackHandler {
	return &FeedbackHandlerImpl{
		feedbackService: feedbackService,
	}
}

// CreateFeedback godoc
// @Summary Create feedback
// @Description Create feedback for a completed booking
// @Tags feedback
// @Accept json
// @Produce json
// @Param request body model.CreateFeedbackRequest true "Feedback creation request"
// @Success 201 {object} ginext.Response{data=model.FeedbackResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/feedback [post]
func (h *FeedbackHandlerImpl) CreateFeedback(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateFeedbackRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	feedback, err := h.feedbackService.CreateFeedback(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create feedback")
		return nil, err
	}

	return ginext.NewCreatedResponse(feedback), nil
}

// GetBookingFeedback godoc
// @Summary Get booking feedback
// @Description Get feedback for a specific booking
// @Tags feedback
// @Produce json
// @Param booking_id path string true "Booking ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.FeedbackResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/feedback/booking/{booking_id} [get]
func (h *FeedbackHandlerImpl) GetBookingFeedback(r *ginext.Request) (*ginext.Response, error) {
	bookingIDStr := r.GinCtx.Param("booking_id")

	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	feedback, err := h.feedbackService.GetBookingFeedback(r.Context(), bookingID)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("failed to get booking feedback")
		return nil, err
	}

	return ginext.NewSuccessResponse(feedback), nil
}

// GetTripFeedbacks godoc
// @Summary Get trip feedbacks
// @Description Get all feedbacks for a specific trip with pagination
// @Tags feedback
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Success 200 {object} ginext.Response{data=model.PaginatedFeedbackResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/feedback/trip/{trip_id} [get]
func (h *FeedbackHandlerImpl) GetTripFeedbacks(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("invalid trip id")
		return nil, ginext.NewBadRequestError("invalid trip id")
	}

	var req model.GetTripFeedbacksRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind query parameters")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	req.Normalize()

	feedbacks, total, err := h.feedbackService.GetTripFeedbacks(r.Context(), tripID, req.Page, req.PageSize)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("failed to get trip feedbacks")
		return nil, err
	}

	return ginext.NewPaginatedResponse(feedbacks, req.Page, req.PageSize, total), nil
}
