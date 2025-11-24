package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockTokenBlacklistManager is a mock implementation of service.TokenBlacklistManager
type MockTokenBlacklistManager struct {
	mock.Mock
}

func (m *MockTokenBlacklistManager) BlacklistToken(ctx context.Context, token string) bool {
	args := m.Called(ctx, token)
	return args.Bool(0)
}

func (m *MockTokenBlacklistManager) IsTokenBlacklisted(ctx context.Context, token string) bool {
	args := m.Called(ctx, token)
	return args.Bool(0)
}

func (m *MockTokenBlacklistManager) BlacklistUserTokens(ctx context.Context, userID uuid.UUID) bool {
	args := m.Called(ctx, userID)
	return args.Bool(0)
}

func (m *MockTokenBlacklistManager) IsUserTokensBlacklisted(ctx context.Context, userID uuid.UUID, tokenIssuedAt int64) bool {
	args := m.Called(ctx, userID, tokenIssuedAt)
	return args.Bool(0)
}
