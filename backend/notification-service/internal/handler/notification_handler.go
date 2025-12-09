package handler

import (
	"github.com/rs/zerolog/log"

	"bus-booking/notification-service/internal/model"
	"bus-booking/notification-service/internal/service"
	"bus-booking/shared/ginext"
)

type NotificationHandler interface {
	Send(r *ginext.Request) (*ginext.Response, error)
}

type NotificationHandlerImpl struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) NotificationHandler {
	return &NotificationHandlerImpl{
		service: service,
	}
}

// Send godoc
// @Summary Send a generic notification
// @Description Send a notification based on type
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body model.GenericNotificationRequest true "Notification request"
// @Success 200 {object} ginext.Response "Notification sent successfully"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /notifications [post]
func (h *NotificationHandlerImpl) Send(r *ginext.Request) (*ginext.Response, error) {
	var req model.GenericNotificationRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.service.SendNotification(r.GinCtx.Request.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Failed to send notification")
		return nil, ginext.NewInternalServerError(err.Error())
	}

	return ginext.NewSuccessResponse("Notification sent successfully"), nil
}
