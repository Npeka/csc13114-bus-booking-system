package appinit

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog/log"
)

// InitFirebase initializes Firebase Auth client
func InitFirebase() (*auth.Client, error) {
	// Initialize Firebase app and get auth client
	// Note: Firebase SDK will automatically use GOOGLE_APPLICATION_CREDENTIALS env var
	// or Application Default Credentials
	firebaseApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase Auth client: %w", err)
	}

	log.Info().Msg("Firebase Auth client initialized successfully")
	return authClient, nil
}
