package service

import (
	"context"
)

type NotificationService interface {
	CreateNotification(ctx context.Context) error
}

type NotificationServiceImpl struct {
}

func NewNotificationService() NotificationService {
	return &NotificationServiceImpl{}
}

func (n *NotificationServiceImpl) CreateNotification(ctx context.Context) error {
	panic("unimplemented")
}
