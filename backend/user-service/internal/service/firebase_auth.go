package service

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type FirebaseAuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}
