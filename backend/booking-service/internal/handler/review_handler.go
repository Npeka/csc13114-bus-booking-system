package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ReviewHandler interface {
	CreateReview(r *ginext.Request) (*ginext.Response, error)
	GetReviewByBooking(r *ginext.Request) (*ginext.Response, error)
	GetUserReviews(r *ginext.Request) (*ginext.Response, error)
	GetTripReviews(r *ginext.Request) (*ginext.Response, error)
	GetTripReviewSummary(r *ginext.Request) (*ginext.Response, error)
	UpdateReview(r *ginext.Request) (*ginext.Response, error)
	DeleteReview(r *ginext.Request) (*ginext.Response, error)
	ModerateReview(r *ginext.Request) (*ginext.Response, error)
}

type ReviewHandlerImpl struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) ReviewHandler {
	return &ReviewHandlerImpl{
		reviewService: reviewService,
	}
}

// CreateReview godoc
// @Summary Create review
// @Description Create review for a confirmed booking
// @Tags reviews
// @Accept json
// @Produce json
// @Param booking_id path string true "Booking ID" format(uuid)
// @Param request body model.CreateReviewRequest true "Review creation request"
// @Success 201 {object} ginext.Response{data=model.ReviewResponse}
// @Failure 400 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/{id}/review [post]
// @Security BearerAuth
func (h *ReviewHandlerImpl) CreateReview(r *ginext.Request) (*ginext.Response, error) {
	bookingIDStr := r.GinCtx.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := r.GinCtx.Get("user_id")
	if !exists {
		return nil, ginext.NewUnauthorizedError("user not authenticated")
	}

	var req model.CreateReviewRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// Override booking_id from path
	req.BookingID = bookingID

	review, err := h.reviewService.CreateReview(r.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create review")
		return nil, err
	}

	return ginext.NewCreatedResponse(review), nil
}

// GetReviewByBooking godoc
// @Summary Get booking review
// @Description Get review for a specific booking
// @Tags reviews
// @Produce json
// @Param id path string true "Booking ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.ReviewResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/{id}/review [get]
func (h *ReviewHandlerImpl) GetReviewByBooking(r *ginext.Request) (*ginext.Response, error) {
	bookingIDStr := r.GinCtx.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	review, err := h.reviewService.GetReviewByBooking(r.Context(), bookingID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get review")
		return nil, err
	}

	return ginext.NewSuccessResponse(review), nil
}

// GetUserReviews godoc
// @Summary Get user reviews
// @Description Get all reviews by a user
// @Tags reviews
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query string false "Review status"
// @Success 200 {object} ginext.Response{data=[]model.ReviewResponse,meta=ginext.MetaData}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/users/{user_id}/reviews [get]
func (h *ReviewHandlerImpl) GetUserReviews(r *ginext.Request) (*ginext.Response, error) {
	userIDStr := r.GinCtx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user id")
	}

	var req model.GetUserReviewsRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	req.UserID = userID
	req.Normalize()

	reviews, total, err := h.reviewService.GetUserReviews(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user reviews")
		return nil, err
	}

	return ginext.NewPaginatedResponse(reviews, req.Page, req.PageSize, total), nil
}

// GetTripReviews godoc
// @Summary Get trip reviews
// @Description Get all reviews for a trip
// @Tags reviews
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Param trip_instance_id query string false "Trip Instance ID" format(uuid)
// @Param min_rating query int false "Minimum rating" minimum(1) maximum(5)
// @Param status query string false "Review status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} ginext.Response{data=[]model.ReviewResponse,meta=ginext.MetaData}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/trips/{trip_id}/reviews [get]
func (h *ReviewHandlerImpl) GetTripReviews(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip id")
	}

	var req model.GetTripReviewsRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	req.TripID = &tripID
	req.Normalize()

	reviews, total, err := h.reviewService.GetTripReviews(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get trip reviews")
		return nil, err
	}

	return ginext.NewPaginatedResponse(reviews, req.Page, req.PageSize, total), nil
}

// GetTripReviewSummary godoc
// @Summary Get trip review summary
// @Description Get aggregated review statistics for a trip
// @Tags reviews
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.TripReviewSummary}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/trips/{trip_id}/reviews/summary [get]
func (h *ReviewHandlerImpl) GetTripReviewSummary(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip id")
	}

	summary, err := h.reviewService.GetTripReviewSummary(r.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get review summary")
		return nil, err
	}

	return ginext.NewSuccessResponse(summary), nil
}

// UpdateReview godoc
// @Summary Update review
// @Description Update user's own review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path string true "Review ID" format(uuid)
// @Param request body model.UpdateReviewRequest true "Update request"
// @Success 200 {object} ginext.Response{data=model.ReviewResponse}
// @Failure 400 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/reviews/{id} [put]
// @Security BearerAuth
func (h *ReviewHandlerImpl) UpdateReview(r *ginext.Request) (*ginext.Response, error) {
	reviewIDStr := r.GinCtx.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid review id")
	}

	userID, exists := r.GinCtx.Get("user_id")
	if !exists {
		return nil, ginext.NewUnauthorizedError("user not authenticated")
	}

	var req model.UpdateReviewRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	review, err := h.reviewService.UpdateReview(r.Context(), userID.(uuid.UUID), reviewID, &req)
	if err != nil {
		log.Error().Err(err).Msg("failed to update review")
		return nil, err
	}

	return ginext.NewSuccessResponse(review), nil
}

// DeleteReview godoc
// @Summary Delete review
// @Description Delete user's own review
// @Tags reviews
// @Param id path string true "Review ID" format(uuid)
// @Success 204 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/reviews/{id} [delete]
// @Security BearerAuth
func (h *ReviewHandlerImpl) DeleteReview(r *ginext.Request) (*ginext.Response, error) {
	reviewIDStr := r.GinCtx.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid review id")
	}

	userID, exists := r.GinCtx.Get("user_id")
	if !exists {
		return nil, ginext.NewUnauthorizedError("user not authenticated")
	}

	if err := h.reviewService.DeleteReview(r.Context(), userID.(uuid.UUID), reviewID); err != nil {
		log.Error().Err(err).Msg("failed to delete review")
		return nil, err
	}

	return ginext.NewNoContentResponse(), nil
}

// ModerateReview godoc
// @Summary Moderate review (Admin)
// @Description Change review status and add admin notes
// @Tags admin-reviews
// @Accept json
// @Produce json
// @Param id path string true "Review ID" format(uuid)
// @Param request body map[string]string true "Moderation request" example({"status":"hidden","admin_notes":"Inappropriate content"})
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/admin/reviews/{id}/moderate [put]
// @Security BearerAuth
func (h *ReviewHandlerImpl) ModerateReview(r *ginext.Request) (*ginext.Response, error) {
	reviewIDStr := r.GinCtx.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid review id")
	}

	var req struct {
		Status     string `json:"status" binding:"required,oneof=active hidden flagged removed"`
		AdminNotes string `json:"admin_notes"`
	}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	status := model.ReviewStatus(req.Status)
	if err := h.reviewService.ModerateReview(r.Context(), reviewID, status, req.AdminNotes); err != nil {
		log.Error().Err(err).Msg("failed to moderate review")
		return nil, err
	}

	return ginext.NewSuccessResponse(map[string]string{
		"message": "review moderated successfully",
	}), nil
}
