package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockNotificationClient is a mock implementation of client.NotificationClient
type MockNotificationClient struct {
	mock.Mock
}

func (m *MockNotificationClient) Send(ctx context.Context, email, name, otp string) error {
	args := m.Called(ctx, email, name, otp)
	return args.Error(0)
}
