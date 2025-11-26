package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockPasswordResetService is a mock implementation of PasswordResetService
type MockPasswordResetService struct {
	mock.Mock
}

func (m *MockPasswordResetService) GenerateResetToken(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordResetService) ValidateResetToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordResetService) InvalidateResetToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
