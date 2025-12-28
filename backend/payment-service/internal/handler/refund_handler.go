package handler

import (
	"bus-booking/payment-service/internal/model"
	"bus-booking/payment-service/internal/service"
	sharedcontext "bus-booking/shared/context"
	"bus-booking/shared/ginext"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RefundHandler interface {
	Create(r *ginext.Request) (*ginext.Response, error)
	GetByBookingID(r *ginext.Request) (*ginext.Response, error)
	ListRefunds(r *ginext.Request) (*ginext.Response, error)
	UpdateRefundStatus(r *ginext.Request) (*ginext.Response, error)
	ExportRefunds(r *ginext.Request) error
}

type RefundHandlerImpl struct {
	service service.RefundService
}

func NewRefundHandler(service service.RefundService) RefundHandler {
	return &RefundHandlerImpl{
		service: service,
	}
}

// Create godoc
// @Summary Create a refund request
// @Description Create a refund request for a cancelled booking
// @Tags refunds
// @Accept json
// @Produce json
// @Param refund body model.RefundRequest true "Refund request"
// @Success 201 {object} ginext.Response{data=model.RefundResponse}
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/refunds [post]
func (h *RefundHandlerImpl) Create(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	var req model.RefundRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	refund, err := h.service.CreateRefund(r.Context(), &req, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create refund")
		return nil, err
	}

	return ginext.NewCreatedResponse(refund), nil
}

// GetByBookingID godoc
// @Summary Get refund by booking ID
// @Description Get refund information for a specific booking
// @Tags refunds
// @Accept json
// @Produce json
// @Param booking_id path string true "Booking ID"
// @Success 200 {object} ginext.Response{data=model.RefundResponse}
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/refunds/booking/{booking_id} [get]
func (h *RefundHandlerImpl) GetByBookingID(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)

	bookingIDStr := r.GinCtx.Param("booking_id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid booking ID")
	}

	refund, err := h.service.GetRefundByBookingID(r.Context(), bookingID, userID)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("Failed to get refund")
		return nil, err
	}

	return ginext.NewSuccessResponse(refund), nil
}

// ListRefunds godoc
// @Summary List refunds (Admin)
// @Description List refund transactions with filters
// @Tags admin
// @Accept json
// @Produce json
// @Param status query string false "Refund status" Enums(PENDING, PROCESSING, COMPLETED, REJECTED)
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/refunds [get]
func (h *RefundHandlerImpl) ListRefunds(r *ginext.Request) (*ginext.Response, error) {
	var query model.RefundListQuery
	if err := r.GinCtx.ShouldBindQuery(&query); err != nil {
		log.Debug().Err(err).Msg("Query binding failed")
		return nil, ginext.NewBadRequestError("Invalid query parameters")
	}

	// Normalize defaults
	query.Normalize()

	refunds, total, err := h.service.ListRefunds(r.Context(), &query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list refunds")
		return nil, err
	}

	return ginext.NewPaginatedResponse(refunds, query.Page, query.PageSize, total), nil
}

// UpdateRefundStatus godoc
// @Summary Update refund status (Admin)
// @Description Update the status of a refund transaction
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Refund Transaction ID"
// @Param request body model.UpdateRefundStatusRequest true "Status update request"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/refunds/{id} [put]
func (h *RefundHandlerImpl) UpdateRefundStatus(r *ginext.Request) (*ginext.Response, error) {
	adminID := sharedcontext.GetUserID(r.GinCtx)

	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid refund ID")
	}

	var req model.UpdateRefundStatusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.service.UpdateRefundStatus(r.Context(), id, req.Status, adminID); err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Failed to update refund status")
		return nil, err
	}

	return ginext.NewSuccessResponse("Refund status updated successfully"), nil
}

// ExportRefunds godoc
// @Summary Export refunds to Excel (Admin)
// @Description Export selected refund transactions to an Excel file
// @Tags admin
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param request body model.ExportRefundsRequest true "Export request"
// @Success 200 {file} binary "Excel file"
// @Failure 400 {object} ginext.Response
// @Failure 401 {object} ginext.Response
// @Failure 403 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/refunds/export [post]
func (h *RefundHandlerImpl) ExportRefunds(r *ginext.Request) error {
	var req model.ExportRefundsRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		resp := map[string]interface{}{
			"success": false,
			"message": "Invalid request data",
		}
		r.GinCtx.JSON(400, resp)
		return nil
	}

	excelData, err := h.service.ExportRefundsToExcel(r.Context(), req.RefundIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to export refunds")
		// Return error as JSON
		resp := map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		r.GinCtx.JSON(500, resp)
		return nil
	}

	// Set headers for file download
	filename := fmt.Sprintf("refunds_%s.xlsx", time.Now().Format("20060102_150405"))
	r.GinCtx.Header("Content-Description", "File Transfer")
	r.GinCtx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	r.GinCtx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	r.GinCtx.Header("Content-Transfer-Encoding", "binary")
	r.GinCtx.Header("Expires", "0")
	r.GinCtx.Header("Cache-Control", "must-revalidate")
	r.GinCtx.Header("Pragma", "public")

	// Write file
	r.GinCtx.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)

	return nil
}
