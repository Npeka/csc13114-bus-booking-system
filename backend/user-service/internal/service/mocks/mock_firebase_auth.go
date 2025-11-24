package mocks

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/mock"
)

// MockFirebaseAuthClient is a mock implementation of Firebase Auth Client
type MockFirebaseAuthClient struct {
	mock.Mock
}

func (m *MockFirebaseAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	args := m.Called(ctx, idToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Token), args.Error(1)
}
