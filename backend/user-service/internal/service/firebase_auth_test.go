package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: FirebaseAuth is just a thin wrapper around Firebase client
// These are minimal tests since actual Firebase verification logic is in the Firebase SDK

func TestNewFirebaseAuth(t *testing.T) {
	// Pass nil client for testing (real client requires Firebase credentials)
	firebaseAuth := NewFirebaseAuth(nil)

	assert.NotNil(t, firebaseAuth)
	assert.IsType(t, &FirebaseAuthImpl{}, firebaseAuth)
}

func TestFirebaseAuthImpl_VerifyIDToken_NilClient(t *testing.T) {
	// Test with nil client - should panic or return error
	firebaseAuth := &FirebaseAuthImpl{client: nil}
	ctx := context.Background()

	// This will panic with nil client, so we test for panic
	assert.Panics(t, func() {
		_, _ = firebaseAuth.VerifyIDToken(ctx, "fake-token")
	})
}

// For real Firebase testing, you would need:
// 1. Firebase credentials
// 2. Valid ID tokens
// 3. Integration tests
// Since FirebaseAuth is just a passthrough to Firebase SDK,
// comprehensive unit tests are limited without mocking Firebase client
