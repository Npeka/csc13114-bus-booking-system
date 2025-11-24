package mocks

import (
	"bus-booking/user-service/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockJWTManager is a mock implementation of utils.JWTManager
type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
	args := m.Called(userID, email, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error) {
	args := m.Called(userID, email, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) ValidateAccessToken(tokenString string) (*utils.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}

func (m *MockJWTManager) ValidateRefreshToken(tokenString string) (*utils.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}
