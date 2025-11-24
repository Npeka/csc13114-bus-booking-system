package service

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

// FirebaseAuthClient is an interface that wraps the Firebase Auth Client methods we use
type FirebaseAuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}
