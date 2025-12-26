package service

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type FirebaseAuth interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type FirebaseAuthImpl struct {
	client *auth.Client
}

func NewFirebaseAuth(client *auth.Client) FirebaseAuth {
	return &FirebaseAuthImpl{client: client}
}

func (f *FirebaseAuthImpl) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return f.client.VerifyIDToken(ctx, idToken)
}
