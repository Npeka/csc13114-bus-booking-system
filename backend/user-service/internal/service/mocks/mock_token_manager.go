package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockTokenManager is a mock implementation of service.TokenBlacklistManager
type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) Blacklist(ctx context.Context, token string) bool {
	args := m.Called(ctx, token)
	return args.Bool(0)
}

func (m *MockTokenManager) IsBlacklisted(ctx context.Context, token string) bool {
	args := m.Called(ctx, token)
	return args.Bool(0)
}

func (m *MockTokenManager) BlacklistUserTokens(ctx context.Context, userID uuid.UUID) bool {
	args := m.Called(ctx, userID)
	return args.Bool(0)
}

func (m *MockTokenManager) IsUserTokensBlacklisted(ctx context.Context, userID uuid.UUID, tokenIssuedAt int64) bool {
	args := m.Called(ctx, userID, tokenIssuedAt)
	return args.Bool(0)
}
