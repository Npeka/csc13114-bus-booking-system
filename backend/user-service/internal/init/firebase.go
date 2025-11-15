package appinit

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog/log"
)

func InitFirebase() (*auth.Client, error) {
	firebaseApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Warn().Err(err).Msg("Firebase app initialization failed - running without Firebase Auth")
		return nil, nil // Return nil client without error for development
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("Firebase Auth client initialization failed - running without Firebase Auth")
		return nil, nil // Return nil client without error for development
	}

	log.Info().Msg("Firebase Auth client initialized successfully")
	return authClient, nil
}
