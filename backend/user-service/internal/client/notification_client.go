package client

import (
	"bus-booking/shared/client"
	"bus-booking/user-service/internal/model/notification"
	"context"
	"fmt"
)

type NotificationClient interface {
	Send(ctx context.Context, email, name, otp string) error
}

type NotificationClientImpl struct {
	http client.HTTPClient
}

func NewNotificationClient(serviceName, baseURL string) NotificationClient {
	httpClient := client.NewHTTPClient(&client.Config{
		ServiceName: serviceName,
		BaseURL:     baseURL,
	})

	return &NotificationClientImpl{
		http: httpClient,
	}
}

func (c *NotificationClientImpl) Send(ctx context.Context, email, name, otp string) error {
	req := &notification.GenericNotificationRequest{
		Type: "OTP",
		Payload: map[string]interface{}{
			"email": email,
			"name":  name,
			"otp":   otp,
		},
	}

	_, err := c.http.Post(ctx, "/api/v1/notifications", req, nil)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}
