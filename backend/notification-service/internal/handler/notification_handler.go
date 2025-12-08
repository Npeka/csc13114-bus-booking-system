package handler

import (
	"github.com/rs/zerolog/log"

	"bus-booking/notification-service/internal/service"
	"bus-booking/shared/ginext"
)

type NotificationHandler interface {
	CreateNotification(r *ginext.Request) (*ginext.Response, error)
}

type NotificationHandlerImpl struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) NotificationHandler {
	return &NotificationHandlerImpl{
		service: service,
	}
}

// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a new notification
// @Tags notifications
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notifications [post]
func (h *NotificationHandlerImpl) CreateNotification(r *ginext.Request) (*ginext.Response, error) {
	if err := h.service.CreateNotification(r.GinCtx.Request.Context()); err != nil {
		log.Error().Err(err).Msg("Failed to create notification")
		return nil, ginext.NewInternalServerError("Failed to create notification")
	}

	return ginext.NewSuccessResponse("Notification created successfully"), nil
}
